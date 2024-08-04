package websocket

import (
	"log"
	"net/http"

	gorillaWebsocket "github.com/gorilla/websocket"
)

var upgrader = gorillaWebsocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func HandleConnections(
	manager *HubManager,
	writer http.ResponseWriter,
	request *http.Request,
) {
	log.Println("HandleConnections called")

	matchId := request.URL.Query().Get("matchId")
	if matchId == "" {
		log.Println("Match ID is required")
		return
	}

	log.Println("Match ID:", matchId)

	websocket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
		return
	}
	defer websocket.Close()

	log.Println("WebSocket connection upgraded")

	client := &Client{
		connection: websocket,
		send:       make(chan []byte, 256),
	}

	manager.JoinHub(matchId, client)
	client.hub = manager.GetHub(matchId)

	log.Println("Client joined hub")

	go client.writePump()
	client.readPump()

	log.Println("Client left hub")
	manager.LeaveHub(matchId, client)
}
