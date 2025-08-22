package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

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

	r.Post("/media", triggerMedia())
	r.Post("/webrtc/offer", webrtcOfferHandler)

	return r
}

func serveWs(hub *ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &ChatClient{hub: hub, conn: conn, message: make(chan IncomingMessage, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func triggerMedia() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload TriggerMediaPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.ResponseError(w, "Failed to decode body payload", 400, err)
			return
		}

		go func() {
			// Download S3 file locally
			localPath := "/tmp/" + payload.Key
			path, err := utils.DownloadS3Object(&payload.Bucket, &payload.Key, localPath)

			if err != nil {
				log.Println("failed to download S3 object:", err)
				return
			}

			// Start broadcaster
			ctx := context.Background()

			if err := startBroadcaster(ctx, path); err != nil {
				log.Printf("startBroadcaster error: %v", err)
			}

			// Remove temp file when done
			_ = os.Remove(path)
		}()

		resp, _ := json.Marshal(map[string]any {
			"message": "Trigger media successfully",
		})

		w.Write(resp)
	})
}
