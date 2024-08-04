package websocket

import (
	"encoding/json"
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

	go client.writePump()
	_, message, err := websocket.ReadMessage()
	if err != nil {
		log.Println("Error reading message", err)
		return
	}

	var unmarshaledMessage Message
	if err := json.Unmarshal(message, &unmarshaledMessage); err != nil {
		log.Println("Error unmarshaling message", err)
		websocket.Close()
		return
	}

	if unmarshaledMessage.Type != Initial {
		log.Println("Initial message type required")
		websocket.Close()
		return
	}

	payload, ok := unmarshaledMessage.Payload.(map[string]interface{})
	if !ok {
		log.Println("Invalid payload")
		websocket.Close()
		return
	}

	playerId, ok := payload["playerId"].(string)
	if !ok || playerId == "" {
		log.Println("Player ID is required")
		websocket.Close()
		return
	}

	client.playerId = playerId
	hub := manager.GetHub(matchId)

	if hub.playerIds[playerId] {
		client.connection.WriteMessage(gorillaWebsocket.TextMessage, []byte("Player ID already exists"))
		client.connection.Close()
		return
	}

	manager.JoinHub(matchId, client)

	log.Println("Client joined hub")

	client.readPump()

	log.Println("Client left hub")
	manager.LeaveHub(matchId, client)
}
