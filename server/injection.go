package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
