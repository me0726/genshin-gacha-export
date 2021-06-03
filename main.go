package main

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func main() {
	listen, err2 := net.Listen("tcp", ":0")
	if err2 != nil {
		panic(err2)
	}
	port := listen.Addr().(*net.TCPAddr).Port
	err3 := listen.Close()
	if err3 != nil {
		panic(err3)
	}
	http.HandleFunc("/static/data.xlsx", IndexHandler)
	http.Handle("/", http.FileServer(getFileSystem()))
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
		log.Printf("start success in http://localhost:%d", port)
	}()
	err := exec.Command(`cmd`, `/c`, `start`, fmt.Sprintf("http://localhost:%d", port)).Start()
	if err != nil {
		panic(err)
	}
	select {}
}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	defer func() {
		err := recover()
		if err != nil {

			switch err.(type) {
			case error:
				log.Print(err)
				_, err := w.Write(bytes.NewBufferString(err.(error).Error()).Bytes())
				if err != nil {
					log.Print(err)
				}
			default:
				log.Print(err)
			}
		}
	}()
	dir, e := os.UserHomeDir()
	if e != nil {
		panic(e)
	}
	var target = dir + "\\AppData\\LocalLow\\miHoYo\\原神\\output_log.txt"
	log.Print("使用文件: " + target)
	readFile, err := ioutil.ReadFile(target)
	if err != nil {
		panic(err)
	}
	webUrl := parseWebUrl(string(readFile))
	log.Print("获取到webUrl: " + webUrl)
	jsonUrl := parseJsonUrl(webUrl)
	log.Print("解析JsonUrl: " + jsonUrl)
	file := writeJson(jsonUrl)
	err = file.Write(w)
	if err != nil {
		panic(err)
	}
}

//go:embed ui/build
var webRoot embed.FS

func getFileSystem() http.FileSystem {
	log.Print("using embed mode")
	embedFs, err := fs.Sub(webRoot, "ui/build")
	if err != nil {
		panic(err)
	}
	return http.FS(embedFs)
}
func writeJson(jsonUrl string) *excelize.File {
	excel := excelize.NewFile()
	writeToExcel("新手祈愿", jsonUrl, 100, excel)
	writeToExcel("常驻祈愿", jsonUrl, 200, excel)
	writeToExcel("角色活动祈愿", jsonUrl, 301, excel)
	writeToExcel("武器活动祈愿", jsonUrl, 302, excel)
	excel.DeleteSheet("Sheet1")
	return excel
}

func writeToExcel(sheetName string, jsonUrl string, gachaType int, excel *excelize.File) {
	arr := RealDataList{}
	for i := 0; i < 100; i++ {
		log.Printf("正在解析 <%s> 第 <%d> 页。每页大小：<%d>", sheetName, i+1, 20)
		url := jsonUrl + fmt.Sprintf("&size=20&gacha_type=%d&page=%d", gachaType, i+1)
		if len(arr) == 0 {
			url += "&end_id=0"
		} else {
			url += "&end_id=" + arr[len(arr)-1].Id
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		data := Data{}
		err = json.Unmarshal(content, &data)
		if err != nil {
			panic(err)
		}
		if data.Retcode != 0 {
			panic(errors.New(data.Message))
		}
		list := data.Data.RealDataList
		log.Printf("解析    <%s> 第 <%d> 页完毕. 当前页条数：<%d>", sheetName, i+1, len(list))
		if len(list) == 0 {
			break
		}
		if len(arr) != 0 && list[len(list)-1].Id == arr[len(arr)-1].Id {
			break
		}
		arr = append(arr, list...)
	}
	sort.Stable(arr)
	log.Printf("%s: %d 条.", sheetName, arr.Len())
	excel.SetActiveSheet(excel.NewSheet(sheetName))
	setCellValue(excel, sheetName, "A1", "时间")
	setCellValue(excel, sheetName, "B1", "名称")
	setCellValue(excel, sheetName, "C1", "类别")
	setCellValue(excel, sheetName, "D1", "星级")
	setCellValue(excel, sheetName, "E1", "总次数")
	setCellValue(excel, sheetName, "F1", "保底内")
	count := 1
	for i, data := range arr {
		setCellValue(excel, sheetName, "A"+strconv.Itoa(i+2), data.Time)
		setCellValue(excel, sheetName, "B"+strconv.Itoa(i+2), data.Name)
		setCellValue(excel, sheetName, "C"+strconv.Itoa(i+2), data.ItemType)
		setCellValue(excel, sheetName, "D"+strconv.Itoa(i+2), data.RankType)
		//总次数
		setCellValue(excel, sheetName, "E"+strconv.Itoa(i+2), i+1)
		//保底内
		setCellValue(excel, sheetName, "F"+strconv.Itoa(i+2), count)
		if data.RankType == "5" {
			count = 1
		} else {
			count++
		}
	}
}

func setCellValue(excel *excelize.File, sheetName string, cell string, value interface{}) {
	err := excel.SetCellValue(sheetName, cell, value)
	if err != nil {
		panic(err)
	}
}

func parseWebUrl(content string) string {
	keyword := "OnGetWebViewPageFinish:"
	lines := strings.Split(content, "\n")
	for _, v := range lines {
		if strings.Contains(v, keyword) {
			return v[len(keyword):]
		}
	}
	panic(errors.New("解析失败找不到对应的url"))
}

const BaseUrl = "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"

func parseJsonUrl(webUrl string) string {
	start := strings.Index(webUrl, "?")
	end := strings.Index(webUrl, "#/log")
	if end-start <= 0 {
		panic(errors.New("解析jsonUrl失败"))
	}
	return BaseUrl + webUrl[start:end]
}

type Data struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Page         string     `json:"page"`
		Size         string     `json:"size"`
		Total        string     `json:"total"`
		RealDataList []RealData `json:"list"`
		Region       string     `json:"region"`
	} `json:"data"`
}
type RealData struct {
	Uid       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	Id        string `json:"id"`
}

type RealDataList []RealData

func (x RealDataList) Len() int      { return len(x) }
func (x RealDataList) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x RealDataList) Less(i, j int) bool {
	return x[i].Time < x[j].Time
}
