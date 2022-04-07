package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
)

type handler struct {
	service project.Service
}

func NewHandler(s project.Service) (*handler, error) {
	if s == nil {
		return nil, errors.New("project: arguments for NewHandler can't be nil")
	}
	return &handler{
		service: s,
	}, nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		rw.Write([]byte("get"))
	case "POST":
		h.post(rw, req)
	case "DELETE":
		rw.Write([]byte("get"))
	case "PUT":
		rw.Write([]byte("get"))
	case "PATCH":
		rw.Write([]byte("get"))
	default:
		rw.Write([]byte("no seponse"))
	}
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		log.Fatal(h)
	}
	defer r.Body.Close()
	f, headers, err := r.FormFile("project")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	name := r.FormValue("name")
	basePath := r.FormValue("basePath")
	isSPA := r.FormValue("isSPA")
	if name == "" || basePath == "" || isSPA == "" {
		w.Write([]byte("incorrect input"))
		return
	}
	c := false
	switch isSPA {
	case "true":
		c = true
	case "false":
		break
	default:
		w.Write([]byte("incorrect input"))
		return
	}
	log.Println(headers)

	p, err := h.service.Deploy(context.Background(), f, name, basePath, c)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&p); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}
