package rest

import (
	"context"
	"errors"
	"log"
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
	log.Println(p)
	for _, v := range p {
		m[v.BasePath] = v
	}
	return &handler{
		service:  s,
		projects: m,
	}, nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	switch req.Method {
	case "POST":
		h.post(rw, req)
		return
	case "GET":
		h.get(rw, req)
		return
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
	s, err := h.parseURL(r.URL.Path)
	if err != nil {
		return
	}
	v, ok := h.projects[s[0]]
	if !ok {
		return
	}
	log.Println(s)
	http.ServeFile(w, r, v.RealPath+"/"+s[1])
}

func (h *handler) parseURL(s string) ([2]string, error) {
	if len(s) <= 1 {
		return [2]string{}, errors.New("incorrect input")
	}
	res := [2]string{}
	basepath := []rune{}
	index := 0
	for i, v := range s[1:] {
		if v == '/' {
			index = i
			break
		}
		basepath = append(basepath, v)
	}
	filepath := s[index+1:]
	res[0] = string(basepath)
	res[1] = filepath
	return res, nil
}
