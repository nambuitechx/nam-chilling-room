package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nambuitechx/nam-chilling-room-server/users"
)

type ChatClient struct {
	conn	*websocket.Conn
	send	chan IncomingMessage
	hub		*ChatHub
}

func (c *ChatClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()

		if err != nil {
			break
		}

		var incomingMessage IncomingMessage

		if err := json.Unmarshal(msg, &incomingMessage); err != nil {
			log.Printf("invalid json from client: %v", err)
            continue
		}

		c.hub.broadcast <- incomingMessage
	}
}

func (c *ChatClient) writePump() {
	for msg := range c.send {
		claims, err := users.ValidateTokenString(msg.TokenString)

		if err != nil {
			log.Printf("failed to validate token string: %v", err)
            continue
		}

		responseMessage := map[string]any {
			"tokenString": msg.TokenString,
			"username": claims.Username,
			"content": msg.Content,
		}

		data, err := json.Marshal(responseMessage)

		if err != nil {
			log.Printf("invalid json to send to client: %v", err)
            continue
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
