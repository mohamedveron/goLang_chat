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

type Client struct {
	conn *websocket.Conn
	flag  bool
}

var connections map[string]Client

func sendAll(msg Message, recv string){
	
	for k, c := range connections{
		if k == recv{

				if err := c.conn.WriteJSON(msg); err != nil{
				return
			}
		}
		
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

	for{
		var obj Message
		err := conn.ReadJSON(&obj)
		

		if err != nil{
			return
		}

			c := Client{conn, true}
			connections[obj.Sender] = c

		log.Println(connections)
		
		if obj.Sender != obj.Recipient{
			sendAll(obj, obj.Recipient)
		}
	}
}

func main() {

    fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	connections = make(map[string]Client)

	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":54321", nil)
}