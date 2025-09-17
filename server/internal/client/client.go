package client

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"blohub/internal/util"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connected to the hub.
type Client struct {
	Hub  HubInterface    // The hub to which the client belongs
	Conn *websocket.Conn // The WebSocket connection
	Send chan []byte     // Channel to send messages to the client
}

// HubInterface defines the methods that a hub must implement for client communication.
type HubInterface interface {
	Register() chan *Client
	Unregister() chan *Client
	Broadcast() chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: tighten for production
	},
}

func Upgrader() *websocket.Upgrader {
	return &upgrader
}

// Reads messages from the client's WebSocket connection and sends them to the hub's broadcast channel.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister() <- c
		c.Conn.Close()
	}()

	// Set up connection parameters
	c.Conn.SetReadLimit(util.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(util.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(util.PongWait))
		return nil
	})

	// Main read loop
	for {
		_, message, err := c.Conn.ReadMessage() // Read message from WebSocket
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(message) // Clean up the message
		c.Hub.Broadcast() <- message       // Send message to hub for broadcasting
	}
}

// Writes messages to the WebSocket connection from the client's send channel.
func (c *Client) WritePump() {
	ticker := time.NewTicker(util.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(util.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(util.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
