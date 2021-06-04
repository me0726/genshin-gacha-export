package main

import (
	"bytes"
	"database/sql"
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func startServer(port int, db *sql.DB) error {
	http.HandleFunc("/static/data.xlsx", func(writer http.ResponseWriter, request *http.Request) {
		indexHandler(writer, request, db)
	})
	http.Handle("/", http.FileServer(getFileSystem()))
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(port), nil)
	if err != nil {
		return err
	}
	log.Printf("start success in http://localhost:%d", port)
	return nil
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

func indexHandler(w http.ResponseWriter, _ *http.Request, db *sql.DB) {
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

	for gachaType, gachaName := range _Map {
		err = fetchDataToDb(jsonUrl, gachaType, gachaName, db)
		if err != nil {
			panic(err)
		}
	}
	excel, err := exportToExcel(db)
	if err != nil {
		panic(err)
	}
	err = excel.Write(w)
	if err != nil {
		panic(err)
	}
}
func randomPort() (int, error) {
	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1, err
	}
	port := listen.Addr().(*net.TCPAddr).Port
	err = listen.Close()
	if err != nil {
		return -1, err
	}
	return port, nil
}
