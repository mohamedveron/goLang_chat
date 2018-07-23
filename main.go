package main

import (
	
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var connections map[*websocket.Conn]bool

func sendAll(msg Message){

	for conn := range connections{
			if err := conn.WriteJSON(msg); err != nil{
			delete(connections, conn)
			return 
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)

    if err != nil {
        log.Println(err)
        return
    }
	connections[conn] = true
	for{
		var obj Message
		err := conn.ReadJSON(&obj)
		log.Println(obj)

		if err != nil{
			delete(connections, conn)
			return
		}
		sendAll(obj)
	}
}

func main() {

    fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	connections = make(map[*websocket.Conn]bool)

	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":54321", nil)
}