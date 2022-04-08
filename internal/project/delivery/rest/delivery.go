package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
)

type handler struct {
	service  project.Service
	projects map[string]project.Project
}

func NewHandler(s project.Service) (*handler, error) {
	if s == nil {
		return nil, errors.New("project: arguments for NewHandler can't be nil")
	}
	m := make(map[string]project.Project)
	p, err := s.List(context.Background())
	if err != nil {
		return nil, err
	}
	for _, v := range p {
		m[v.BasePath] = v
	}
	return &handler{
		service:  s,
		projects: m,
	}, nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		h.post(rw, req)
		return
	case "GET":
		h.get(rw, req)
		return
	default:
		respondString(rw, http.StatusBadRequest, "no response")
	}
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseMultipartForm(200000); err != nil {
		respondString(w, http.StatusBadGateway, err.Error())
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["project"]

	name := r.FormValue("name")
	basePath := r.FormValue("basePath")
	isSPA := r.FormValue("isSPA")
	if name == "" || basePath == "" || isSPA == "" {
		respondString(w, http.StatusBadRequest, "incorrect input")
		return
	}
	c := false
	switch isSPA {
	case "true":
		c = true
	case "false":
		break
	default:
		respondString(w, http.StatusBadRequest, "incorrect input")
		return
	}

	p, err := h.service.Deploy(context.Background(), files, name, basePath, c)
	if err != nil {
		respondString(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, &p)
	h.projects[name] = p
	return
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	if !IsFileURL(r.URL.Path) {
		s, err := parseURL(r.URL.Path)
		if err != nil {
			return
		}
		v, ok := h.projects[s[0]]
		if !ok {
			return
		}
		if !v.IsSPA {
			return
		}
		http.ServeFile(w, r, v.RealPath+"/"+"index.html")
		return
	}
	s, err := parseURL(r.URL.Path)
	if err != nil {
		return
	}
	v, ok := h.projects[s[0]]
	if !ok {
		return
	}
	http.ServeFile(w, r, v.RealPath+"/"+s[1])
}
