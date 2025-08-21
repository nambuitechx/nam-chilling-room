package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/nambuitechx/nam-chilling-room-server/chat"
	"github.com/nambuitechx/nam-chilling-room-server/configs"
	"github.com/nambuitechx/nam-chilling-room-server/users"
)

func NewRouter() *chi.Mux {
	// Setup
	localEnv := configs.NewLocalEnv()
	db := configs.RunMigration(localEnv)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(map[string]any{
			"message": "healthy",
		})

		w.Write(resp)
	})

	// Users
	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(userRepository)
	userRouter := users.NewUserRouter(userService)

	r.Mount("/users", userRouter)

	// Chat
	chatRouter := chat.NewChatRouter()

	r.Mount("/chat", chatRouter)

	return r
}
