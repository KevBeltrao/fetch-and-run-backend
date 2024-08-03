package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func HandleConnections(hub *Hub, writer http.ResponseWriter, request *http.Request) {
	websocket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
	}
	defer websocket.Close()

	client := &Client{
		hub:        hub,
		connection: websocket,
		send:       make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	client.readPump()
}
