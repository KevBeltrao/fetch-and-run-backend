package websocket

import (
	"encoding/json"
	"log"

	gorillaWebsocket "github.com/gorilla/websocket"
)

type PlayerState struct {
	X int
	Y int
}

type GameState struct {
	Players map[string]PlayerState
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	state      GameState
	playerIds  map[string]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		state: GameState{
			Players: make(map[string]PlayerState),
		},
		playerIds: make(map[string]bool),
	}
}

func (hub *Hub) handlePlayerMove(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("error: invalid payload")
		return
	}

	playerId, ok := data["playerId"].(string)
	if !ok {
		log.Printf("error: invalid playerId")
		return
	}

	x, ok := data["x"].(float64)
	if !ok {
		log.Printf("error: invalid x")
		return
	}

	y, ok := data["y"].(float64)
	if !ok {
		log.Printf("error: invalid y")
		return
	}

	hub.state.Players[playerId] = PlayerState{
		X: int(x),
		Y: int(y),
	}

	updatedState, _ := json.Marshal(Message{
		Type:    "updateState",
		Payload: hub.state,
	})

	log.Printf("Broadcasting updated state: %s", string(updatedState))

	for client := range hub.clients {
		select {
		case client.send <- updatedState:
			log.Printf("Message sent to client: %v", client)
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}

func (hub *Hub) handleMessage(message Message) {
	log.Printf("Handling message: %v", message)
	switch message.Type {
	case "playerMove":
		hub.handlePlayerMove(message.Payload)
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			if hub.playerIds[client.playerId] {
				client.connection.WriteMessage(gorillaWebsocket.TextMessage, []byte(`{"error": "duplicate playerId"}`))
				client.connection.Close()
				continue
			}

			hub.clients[client] = true
			hub.playerIds[client.playerId] = true

			log.Printf("Client registered: %v. Total clients: %d", client, len(hub.clients))
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				delete(hub.playerIds, client.playerId)
				close(client.send)
				log.Printf("Client unregistered: %v. Total clients: %d", client, len(hub.clients))
			}
		case message := <-hub.broadcast:
			var messageUnmarshaled Message
			if err := json.Unmarshal(message, &messageUnmarshaled); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			hub.handleMessage(messageUnmarshaled)
		}
	}
}
