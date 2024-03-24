package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Id      int64
	Type    string
	Success bool
	Raw     []byte
}

func ListenWebsocket(conn *websocket.Conn, ctx context.Context, c chan ChanMsg) {
	for {
		fmt.Println("1")
		bytes, err := ReadMessage(conn, ctx)
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			close(c)
			break
		}

		fmt.Println("2")
		base := BaseMessage{
			// Default to true for messages that do not contain the 'success' field
			Success: true,
		}
		fmt.Println("3")
		_ = json.Unmarshal(bytes, &base)
		fmt.Println("4")
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
		}
		fmt.Println("5")
		chanMsg := ChanMsg{
			Type:    base.Type,
			Id:      base.Id,
			Success: base.Success,
			Raw:     bytes,
		}
		fmt.Println("6")
		c <- chanMsg
	}
}
