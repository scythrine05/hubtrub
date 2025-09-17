package hub

import (
	"net/http"

	"blohub/internal/client"
)

type Hub struct {
	Clients     map[*client.Client]bool // Registered clients
	BroadcastC  chan []byte             // Channel for broadcasting messages to clients
	RegisterC   chan *client.Client     // Channel for registering new clients
	UnregisterC chan *client.Client     // Channel for unregistering clients
}

func NewHub() *Hub {
	return &Hub{
		BroadcastC:  make(chan []byte),             // Channel for broadcasting messages
		RegisterC:   make(chan *client.Client),     // Channel for registering clients
		UnregisterC: make(chan *client.Client),     // Channel for unregistering clients
		Clients:     make(map[*client.Client]bool), // Map to track registered clients
	}
}

func (h *Hub) Register() chan *client.Client   { return h.RegisterC }
func (h *Hub) Unregister() chan *client.Client { return h.UnregisterC }
func (h *Hub) Broadcast() chan []byte          { return h.BroadcastC }

// Run starts the hub's event loop, handling client registration, unregistration, and broadcasting messages.
func (h *Hub) Run() {
	for {
		select {

		// Handle client registration
		case client := <-h.RegisterC:
			h.Clients[client] = true
		// Handle client unregistration
		case client := <-h.UnregisterC:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		// Handle broadcasting messages to all clients
		case message := <-h.BroadcastC:
			// Send the message to all registered clients
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func ServeWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := client.Upgrader().Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}

	// Create a new client and register it with the hub
	c := &client.Client{Hub: h, Conn: conn, Send: make(chan []byte, 256)}
	h.Register() <- c

	// Start the client's read and write pumps in separate goroutines
	go c.WritePump()
	go c.ReadPump()
}
