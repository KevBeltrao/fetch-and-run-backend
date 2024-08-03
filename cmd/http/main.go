package main

import (
	"log"
	"net/http"

	"github.com/kevbeltrao/fetch-and-run-backend/internal/server/websocket"
)

const PORT = ":8000"

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/websocket", func(writer http.ResponseWriter, request *http.Request) {
		websocket.HandleConnections(hub, writer, request)
	})

	log.Println("Server started on", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
