package shared

import "net"

type PlayerData struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Speed float64 `json:"speed"`
}

type Client struct {
	ID   string
	Conn net.Conn
	Out  chan string
}
