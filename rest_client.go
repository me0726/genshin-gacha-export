package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

func fetchDataToDb(jsonUrl string, gachaType int, gachaName string, db *sql.DB) error {
	var lastId = "0"
	for i := 0; i < 1000; i++ {
		log.Printf("正在解析 <%s> 第 <%d> 页。每页大小：<%d>", gachaName, i+1, 20)
		url := jsonUrl + fmt.Sprintf("&size=20&gacha_type=%d&page=%d&end_id=%s", gachaType, i+1, lastId)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data := Data{}
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		if data.Retcode != 0 {
			return errors.New(data.Message)
		}
		list := data.Data.RealDataList
		log.Printf("解析    <%s> 第 <%d> 页完毕. 当前页条数：<%d>", gachaName, i+1, len(list))
		if len(list) == 0 {
			break
		}
		for index, realData := range list {
			exists, err := existsById(realData.Id, db)
			if err != nil {
				return err
			}
			if !exists {
				err := insert(db, realData)
				if err != nil {
					return err
				}
			}
			if index == len(list)-1 {
				lastId = realData.Id
			}
		}
	}
	return nil
}
