package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var clients = make(map[net.Conn]bool)
var broadcast = make(chan []byte)
var mutex = sync.Mutex{}

func tcpController(conn net.Conn) {
	defer conn.Close()
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break

		}

		broadcast <- buf[:n]

	}
}

func tcpHandle() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			_, err := client.Write(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()

	}

}

func Tcp12() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
	}
	defer listen.Close()
	go tcpHandle()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("the conn is accepted ")
			continue
		}
		go tcpController(conn)
	}

}
