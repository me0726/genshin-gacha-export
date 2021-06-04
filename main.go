package main

import (
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os/exec"
)

var _Map = map[int]string{
	100: "新手祈愿",
	200: "常驻祈愿",
	301: "角色活动祈愿",
	302: "武器活动祈愿",
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Print(err)
		}
	}()
	db, err := initSqlite()
	if err != nil {
		panic(err)
	}
	port, err := randomPort()
	if err != nil {
		panic(err)
	}
	go func() {
		err := startServer(port, db)
		if err != nil {
			log.Print(err)
		}
	}()
	err = openBrowser(port)
	if err != nil {
		panic(err)
	}
	select {}
}

func openBrowser(port int) error {
	return exec.Command(`cmd`, `/c`, `start`, fmt.Sprintf("http://localhost:%d", port)).Start()
}
