package main

import (
	"emperror.dev/emperror"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"strconv"
)

var port = flag.Int("p", 0, "")
var router = httprouter.New()

func main() {
	flag.Parse()
	if *port == 0 {
		*port = 8080
	}
	l, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	emperror.Panic(err)
	registerRoute()
	http.Serve(l, router)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	origin, host := r.Header.Get("Origin"), r.Host

	fmt.Println("origin: ", origin)
	fmt.Println("host: ", host)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "get,post")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Max-Age", "60")
	w.Header().Set("a", "1")
	w.Header().Set("b", "2")

	fmt.Fprintf(w, `{"result":"success"}`)
}

func registerRoute() {

	router.GET("/apis", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		apis := []string{
			"/simple_request/allowed",
			"/simple_request/rejected",
			"/comp_request/allowed",
			"/comp_request/reject_for_headers",
			"/comp_request/allowed_with_cookie",
			"/comp_request/allowed_with_no_vary",
			"/comp_request/allowed_with_vary",
		}
		data, _ := json.Marshal(apis)
		fmt.Fprintf(writer, string(data))
	})

	// 处理简单请求，
	router.GET("/simple_request/allowed", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// 设置 Access-Control-Allow-Origin 允许任何origin访问
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(writer, "access allowed")
	})
	router.GET("/simple_request/rejected", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// 未设置 Access-Control-Allow-Origin
		fmt.Fprintf(writer, "access rejected")
	})

	router.OPTIONS("/comp_request/allowed", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// 未设置 Access-Control-Allow-Origin
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		writer.Header().Set("Access-Control-Allow-Headers", "*")
		//writer.Header().Set("Access-Control-Max-Age", "60") // 1min
		writer.Header().Set("Vary", "Origin")
	})

	// 没有设置header，被拦截
	router.OPTIONS("/comp_request/reject_for_headers", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// 未设置 Access-Control-Allow-Origin
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		//writer.Header().Set("Access-Control-Max-Age", "60") // 1min
		writer.Header().Set("Vary", "Origin")
	})

	router.OPTIONS("/comp_request/allowed_with_cookie", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		writer.Header().Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//writer.Header().Set("Access-Control-Max-Age", "60") // 1min
		writer.Header().Set("Vary", "Origin")
	})

	router.OPTIONS("/comp_request/allowed_with_no_vary", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		writer.Header().Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
		writer.Header().Set("Access-Control-Max-Age", "120") // 10min
		writer.Header().Set("Vary", "")
	})

	router.OPTIONS("/comp_request/allowed_with_vary", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
		writer.Header().Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Max-Age", "120") // 10min
		writer.Header().Set("Vary", "Origin")
	})

	router.GET("/comp_request/*oper", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))

		switch params.ByName("oper") {
		case "/allowed_with_cookie":
			writer.Header().Set("Access-Control-Allow-Credentials", "true")

		}

		fmt.Fprintf(writer, "comp_request access allowed")
	})

}
