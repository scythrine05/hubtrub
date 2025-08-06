package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/scythrine/gozwet/server/shared"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			break
		}

		var data shared.PlayerData
		if err := json.Unmarshal(msg, &data); err != nil {
			fmt.Println("Invalid data:", err)
			continue
		}

		fmt.Printf("Player %s moved to X:%.2f Y:%.2f Speed:%.2f\n", data.ID, data.X, data.Y, data.Speed)
	}
}
