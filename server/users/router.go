package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nambuitechx/nam-chilling-room-server/utils"
)

func NewUserRouter(userService *UserService) http.Handler {
	r := chi.NewRouter()

	r.Get("/", listUsers(userService))
	r.Get("/{userID}", getUserByID(userService))
	r.Post("/register", createUser(userService))
	r.Post("/login", login(userService))

	return r
}

func listUsers(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		limit, err := getQueryInt(r, "limit", 10)
		
		if err != nil {
			utils.ResponseError(w, "Invalid limit", 400, err)
			return
		}

		offset, err := getQueryInt(r, "offset", 0)
		
		if err != nil {
			utils.ResponseError(w, "Invalid offset", 400, err)
			return
		}

		users, err := s.listUsers(username, limit, offset)

		if err != nil {
			utils.ResponseError(w, "Failed to get all users", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Get all users successfully",
			"data": users,
		})

		w.Write(resp)
	})
}

func getUserByID(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		user, err := s.getUserByID(userID)

		if err != nil {
			utils.ResponseError(w, "Failed to get user by id", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Get user by id successfully",
			"data": *user,
		})

		w.Write(resp)
	})
}

func createUser(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload CreateUserPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.ResponseError(w, "Failed to decode body payload", 400, err)
			return
		}

		user, err := s.createUser(payload.Username, payload.Password)

		if err != nil {
			utils.ResponseError(w, "Failed to create new user", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Create new user successfully",
			"data": *user,
		})

		w.Write(resp)
	})
}

func login(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload LoginPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.ResponseError(w, "Failed to decode body payload", 400, err)
			return
		}

		tokenString, err := s.authenticate(payload.Username, payload.Password)

		if err != nil {
			utils.ResponseError(w, "Invalid username or password", 400, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "login successfully",
			"data": tokenString,
		})

		w.Write(resp)
	})
}

func getQueryInt(r *http.Request, key string, defaultValue int) (int, error) {
    valStr := r.URL.Query().Get(key)

    if valStr == "" {
        return defaultValue, nil
    }

    return strconv.Atoi(valStr)
}
