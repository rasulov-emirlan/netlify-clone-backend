package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func respondString(w http.ResponseWriter, s int, v string) {
	w.WriteHeader(s)
	w.Write([]byte(v))
}

func respondJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func parseParam(s string) (string, error) {
	if len(s) <= 1 {
		return "", errors.New("incorrect input")
	}

	temp := strings.Split(s, "/")
	return temp[len(temp)-1], nil
}
