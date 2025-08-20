package chat

import "github.com/gorilla/websocket"

type ChatClient struct {
	conn *websocket.Conn
	send chan []byte
	hub *ChatHub
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

		message := Message{
			Sender: c,
			Content: msg,
		}

		c.hub.broadcast <- message
	}
}

func (c *ChatClient) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
