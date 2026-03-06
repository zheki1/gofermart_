package helpers

import (
	"encoding/json"
	"net/http"
)

// записываем ошибку в ответ json
func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
