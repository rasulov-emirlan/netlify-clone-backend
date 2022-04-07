package rest

import (
	"encoding/json"
	"net/http"
)

func respondString(w http.ResponseWriter, v string) {
	w.Write([]byte(v))
}

func respondJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
