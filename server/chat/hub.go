package chat

import (
	"log"
	"time"
)

type ChatHub struct {
	clients    map[*ChatClient]bool
	broadcast  chan IncomingMessage
	register   chan *ChatClient
	unregister chan *ChatClient

	// Worker queues (separate channels)
	dbQueue   chan IncomingMessage
}

func newHub(dbWorkers int) *ChatHub {
	h := &ChatHub{
		clients:    make(map[*ChatClient]bool),
		broadcast:  make(chan IncomingMessage),
		register:   make(chan *ChatClient),
		unregister: make(chan *ChatClient),
		dbQueue:    make(chan IncomingMessage),
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
