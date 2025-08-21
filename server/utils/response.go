package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseError(w http.ResponseWriter, message string, statusCode int, err error) {
	resp, _ := json.Marshal(map[string]any {
		"message": message,
		"error": err.Error(),
	})

	w.WriteHeader(statusCode)
	w.Write(resp)
}
