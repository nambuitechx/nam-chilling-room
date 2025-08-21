package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ChatHub struct {
	clients    		map[*ChatClient]bool
	register   		chan *ChatClient
	unregister 		chan *ChatClient

	broadcast  		chan IncomingMessage
	mediaBroadcast 	chan []byte

	// Worker queues (separate channels)
	dbQueue   		chan IncomingMessage
}

func newHub(dbWorkers int) *ChatHub {
	h := &ChatHub{
		clients:    	make(map[*ChatClient]bool),
		register:   	make(chan *ChatClient),
		unregister: 	make(chan *ChatClient),
		broadcast:  	make(chan IncomingMessage),
		mediaBroadcast:	make(chan []byte),
		dbQueue:    	make(chan IncomingMessage),
	}

	// Start separate worker pools
	for i := 0; i < dbWorkers; i++ {
		go h.dbWorker(i)
	}

	return h
}

func (h *ChatHub) run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true

			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
					client.conn.Close()
				}

			case message := <-h.broadcast:
				// Broadcast fast
				for client := range h.clients {
					select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
					}
				}

				// // Send to each specialized worker pool
				// select {
				// 	case h.dbQueue <- message:
				// 	default:
				// 		log.Println("⚠️ DB queue full, dropping message")
				// }
			
			case chunk := <-h.mediaBroadcast:
				// Broadcast media chunks
				for client := range h.clients {
					err := client.conn.WriteMessage(websocket.BinaryMessage, chunk)

					if err != nil {
						log.Println("error sending chunk:", err)
						// client.conn.Close()
						// delete(h.clients, client)
						continue
					}
				}
		}
	}
}

func (h *ChatHub) dbWorker(id int) {
	for msg := range h.dbQueue {
		// Simulate DB insert
		log.Printf("[DB Worker %d] Saving to DB: %s\n", id, msg)
		time.Sleep(500 * time.Millisecond)
	}
}
