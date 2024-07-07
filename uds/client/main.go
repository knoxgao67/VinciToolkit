package main

import (
	"emperror.dev/emperror"
	"fmt"
	"github.com/knoxgao67/VinciToolkit/uds/common"
	"log"
	"net"
)

func main() {
	common.Init(false)
	path := common.SocketPath
	if path == "" {
		log.Fatalln("must set path")
	}
	// Replace "/path/to/socket" with the actual path to the UDS socket
	conn, err := net.Dial("unix", path)
	fmt.Println("conn", conn.LocalAddr(), conn.RemoteAddr())
	emperror.Panic(err)
	defer conn.Close()

	// Send data to server
	_, err = conn.Write([]byte("Hello from client\n"))
	emperror.Panic(err)
	if err != nil {
		fmt.Println("Failed to write:", err)
		return
	}

	// Optionally, receive data from server
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	emperror.Panic(err)

	fmt.Println("Received:", string(response[:n]))
}
