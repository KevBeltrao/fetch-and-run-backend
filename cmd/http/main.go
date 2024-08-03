package main

import (
	"log"
	"net/http"

	"github.com/kevbeltrao/fetch-and-run-backend/internal/server/websocket"
)

func main() {
	http.HandleFunc("/websocket", websocket.HandleConnections)

	log.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
