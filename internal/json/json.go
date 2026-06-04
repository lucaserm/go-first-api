package json

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"data": data,
	})
}

func Read(r *http.Request) int64 {
	return 0
}
