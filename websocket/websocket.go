package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = sync.Mutex{}

func websocketController(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- msg
	}

}

func handlemsg() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Fatal(err)
				client.Close()
				delete(clients, client)

			}

		}
		mutex.Unlock()

	}

}

func Websocket() {
	http.HandleFunc("/ws", websocketController)
	go handlemsg()

	fmt.Println("Starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}

}
