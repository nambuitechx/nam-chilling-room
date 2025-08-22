package chat

import (
	"log"
	"time"
)

type ChatHub struct {
	clients    		map[*ChatClient]bool
	register   		chan *ChatClient
	unregister 		chan *ChatClient

	broadcast  		chan IncomingMessage

	// Worker queues (separate channels)
	dbQueue   		chan IncomingMessage
}

func newHub(dbWorkers int) *ChatHub {
	h := &ChatHub{
		clients:    	make(map[*ChatClient]bool),
		register:   	make(chan *ChatClient),
		unregister: 	make(chan *ChatClient),
		broadcast:  	make(chan IncomingMessage),
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
					close(client.message)
					client.conn.Close()
				}

			case message := <-h.broadcast:
				for client := range h.clients {
					select {
						case client.message <- message:
						default:
							close(client.message)
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
