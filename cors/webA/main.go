package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

//go:embed index.html
var indexhtml embed.FS

var port = flag.Int("p", 0, "")

func main() {
	flag.Parse()
	if *port == 0 {
		*port = 8080
	} // 获取当前工作目录

	// 设置静态文件路径
	fileServer := http.FileServerFS(indexhtml)
	// 创建 HTTP 服务器
	http.Handle("/", fileServer)

	fmt.Println("Starting static file server on http://localhost:8080")
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
