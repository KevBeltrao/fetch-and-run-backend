package websocket

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub        *Hub
	connection *websocket.Conn
	send       chan []byte
}

func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client
		client.connection.Close()
	}()

	client.connection.SetReadLimit(maxMessageSize)
	client.connection.SetReadDeadline(time.Now().Add(pongWait))
	client.connection.SetPongHandler(func(string) error {
		client.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		client.hub.broadcast <- message
	}
}

func (client *Client) writePump() {
	ticket := time.NewTicker(pingPeriod)
	defer func() {
		ticket.Stop()
		client.connection.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			n := len(client.send)

			for i := 0; i < n; i += 1 {
				writer.Write(newline)
				writer.Write(<-client.send)
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticket.C:
			client.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
