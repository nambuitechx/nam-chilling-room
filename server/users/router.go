package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewUserRouter(userService *UserService) http.Handler {
	r := chi.NewRouter()

	r.Get("/", ListUsers(userService))
	r.Get("/{userID}", GetUserByID(userService))
	r.Post("/register", CreateUser(userService))
	r.Post("/login", Login(userService))

	return r
}

func ListUsers(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		limit, err := getQueryInt(r, "limit", 10)
		
		if err != nil {
			responseError(w, "Invalid limit", 400, err)
			return
		}

		offset, err := getQueryInt(r, "offset", 0)
		
		if err != nil {
			responseError(w, "Invalid offset", 400, err)
			return
		}

		users, err := s.ListUsers(username, limit, offset)

		if err != nil {
			responseError(w, "Failed to get all users", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Get all users successfully",
			"data": users,
		})

		w.Write(resp)
	})
}

func GetUserByID(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		user, err := s.GetUserByID(userID)

		if err != nil {
			responseError(w, "Failed to get user by id", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Get user by id successfully",
			"data": *user,
		})

		w.Write(resp)
	})
}

func CreateUser(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload CreateUserPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			responseError(w, "Failed to decode body payload", 400, err)
			return
		}

		user, err := s.CreateUser(payload.Username, payload.Password)

		if err != nil {
			responseError(w, "Failed to create new user", 500, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Create new user successfully",
			"data": *user,
		})

		w.Write(resp)
	})
}

func Login(s *UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload LoginPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			responseError(w, "Failed to decode body payload", 400, err)
			return
		}

		tokenString, err := s.Authenticate(payload.Username, payload.Password)

		if err != nil {
			responseError(w, "Invalid username or password", 400, err)
			return
		}

		resp, _ := json.Marshal(map[string]any {
			"message": "Login successfully",
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

func responseError(w http.ResponseWriter, message string, statusCode int, err error) {
	resp, _ := json.Marshal(map[string]any {
		"message": message,
		"error": err.Error(),
	})

	w.WriteHeader(statusCode)
	w.Write(resp)
}
