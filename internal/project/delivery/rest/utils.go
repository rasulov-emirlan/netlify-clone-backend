package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
)

func respondString(w http.ResponseWriter, s int, v string) {
	w.WriteHeader(s)
	w.Write([]byte(v))
}

func respondJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func IsFileURL(s string) bool {
	if len(filepath.Ext(s)) == 0 {
		return false
	}
	return true
}

func parseURL(s string) ([2]string, error) {
	if len(s) <= 1 {
		return [2]string{}, errors.New("incorrect input")
	}
	res := [2]string{}
	basepath := []rune{}
	index := 0
	for i, v := range s[1:] {
		if v == '/' {
			index = i + 1
			break
		}
		basepath = append(basepath, v)
	}
	filepath := s[index+1:]
	res[0] = string(basepath)
	res[1] = filepath
	return res, nil
}
