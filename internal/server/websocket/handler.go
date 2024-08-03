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

func HandleConnections(writer http.ResponseWriter, request *http.Request) {
	websocket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
	}
	defer websocket.Close()

	for {
		var message map[string]interface{}

		err := websocket.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message", err)
			break
		}

		log.Printf("Received message: %v", message)

		err = websocket.WriteJSON(message)
		if err != nil {
			log.Println("Error writing message", err)
			break
		}
	}
}
