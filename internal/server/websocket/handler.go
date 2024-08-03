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

func HandleConnections(
	manager *HubManager,
	writer http.ResponseWriter,
	request *http.Request,
) {
	matchId := request.URL.Query().Get("matchId")
	if matchId == "" {
		log.Println("Match ID is required")
		return
	}

	websocket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
	}
	defer websocket.Close()

	hub := manager.GetHub(matchId)

	client := &Client{
		hub:        hub,
		connection: websocket,
		send:       make(chan []byte, 256),
	}

	manager.JoinHub(matchId, client)

	go client.writePump()
	client.readPump()

	manager.LeaveHub(matchId, client)
}
