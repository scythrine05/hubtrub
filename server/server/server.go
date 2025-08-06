package server

import (
	"fmt"
	"net"

	"github.com/scythrine/gozwet/server/server"
)

func StartServer() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go server.HandleConnection(conn)
	}
}
