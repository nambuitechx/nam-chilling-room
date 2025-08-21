package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/nambuitechx/nam-chilling-room-server/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewChatRouter() http.Handler {
	hub := newHub(3)
	go hub.run()

	r := chi.NewRouter()

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	r.Post("/media", triggerMedia(hub))

	return r
}

func serveWs(hub *ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &ChatClient{hub: hub, conn: conn, send: make(chan IncomingMessage, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func triggerMedia(hub *ChatHub) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload TriggerMediaPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.ResponseError(w, "Failed to decode body payload", 400, err)
			return
		}

		// Start streaming from S3 in the background
		go func() {
			err := utils.StreamS3Object(payload.Bucket, payload.Key, 6 * 1024, hub.mediaBroadcast)

			if err != nil {
				log.Println("stream error:", err)
			}
		}()

		resp, _ := json.Marshal(map[string]any {
			"message": "Trigger media successfully",
		})

		w.Write(resp)
	})
}
