package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseError(w http.ResponseWriter, message string, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	payload := map[string]any {
		"message": message,
		"error": err.Error(),
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
