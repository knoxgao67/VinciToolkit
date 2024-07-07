package main

import (
	"emperror.dev/emperror"
	"flag"
	"fmt"
	"github.com/knoxgao67/VinciToolkit/uds/common"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var startHttp = flag.Bool("http", false, "start http server")

func main() {
	common.Init(true)
	path := common.SocketPath

	fmt.Println("start server with pid: ", os.Getpid())
	fmt.Println("Listening on", path)
	// 监听Unix Domain Socket
	l, err := net.Listen("unix", path)
	emperror.Panic(err)
	defer l.Close()

	if *startHttp {
		startHttpServer(l)
	} else {
		startUdsServer(l)
	}
}

func startHttpServer(l net.Listener) {
	err := http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		emperror.Panic(err)
		fmt.Println(string(body))
		fmt.Fprintln(w, "Hello from server")
	}))
	emperror.Panic(err)
}

func startUdsServer(l net.Listener) {
	// 接受连接并处理
	go func() {
		for {
			conn, err := l.Accept()
			fmt.Println("conn", conn.LocalAddr(), conn.RemoteAddr())
			//conn.LocalAddr()
			if err != nil {
				panic(err)
			}
			// 处理连接
			go handleConnection(conn)
		}
	}()
	// 处理系统信号以优雅地关闭服务
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Shutting down server")
}

// handleConnection 处理客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r, string(debug.Stack()))
		}
	}()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("EOF", "Close by client")
			return
		}
		emperror.Panic(err)
		if n == 0 {
			fmt.Println("Closed by client")
			return // 客户端关闭连接
		}
		fmt.Println("Received:", string(buf[:n]))
		// 可以在这里添加对数据的处理逻辑
		_, err = conn.Write([]byte("Hello from server"))
		emperror.Panic(err)
	}
}
