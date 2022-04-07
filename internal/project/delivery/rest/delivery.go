package rest

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
)

var projects map[string]project.Project

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
	case "POST":
		h.post(rw, req)
	default:
		rw.Write([]byte("no response"))
	}
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		log.Fatal(h)
	}
	defer r.Body.Close()

	if err := r.ParseMultipartForm(200000); err != nil {
		respondString(w, err.Error())
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["project"]

	name := r.FormValue("name")
	basePath := r.FormValue("basePath")
	isSPA := r.FormValue("isSPA")
	if name == "" || basePath == "" || isSPA == "" {
		respondString(w, "incorrect input")
		return
	}
	c := false
	switch isSPA {
	case "true":
		c = true
	case "false":
		break
	default:
		respondString(w, "incorrect input")
		return
	}

	p, err := h.service.Deploy(context.Background(), files, name, basePath, c)
	if err != nil {
		respondString(w, err.Error())
		return
	}
	respondJSON(w, &p)
	projects[name] = p
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {

}
