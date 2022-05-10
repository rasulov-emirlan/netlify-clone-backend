package rest

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
)

type handler struct {
	service project.Service

	// TODO: use a custom DS instead of sync.Map
	projects sync.Map
}

func NewHandler(s project.Service) (*handler, error) {
	if s == nil {
		return nil, errors.New("project: arguments for NewHandler can't be nil")
	}
	m := sync.Map{}
	p, err := s.List(context.Background())
	if err != nil {
		return nil, err
	}
	for _, v := range p {
		log.Println(v)
		m.Store(v.BasePath, v)
	}
	return &handler{
		service:  s,
		projects: m,
	}, nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.post(rw, req)
		return
	case http.MethodGet:
		h.get(rw, req)
		return
	case http.MethodPatch:
		h.patch(rw, req)
		return
	case http.MethodDelete:
		h.delete(rw, req)
		return
	default:
		respondString(rw, http.StatusBadRequest, "no response")
	}
}

const megabyte int64 = 1048576

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseMultipartForm(megabyte * 20); err != nil {
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
	h.projects.Store(basePath, p)
}

func (h *handler) patch(w http.ResponseWriter, r *http.Request) {
	id, err := parseParam(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := r.ParseMultipartForm(megabyte * 20); err != nil {
		respondString(w, http.StatusBadGateway, err.Error())
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["project"]
	if err := h.service.Redeploy(context.Background(), files, project.ID(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		reqURL     = strings.Split(r.URL.Path, "/")
		forwarding = ""
	)

	// TODO: this code is highly dependent
	// on the path to this handler
	// maybe try changing it somehow
	if len(reqURL) < 2 {
		w.WriteHeader(404)
		return
	}

	v, ok := h.projects.Load(reqURL[2])
	if !ok {
		w.WriteHeader(404)
		return
	}

	p, ok := v.(project.Project)
	if !ok {
		http.Error(w, "could not convert interface{}", http.StatusInternalServerError)
		return
	}
	switch {
	case p.IsSPA && path.Ext(reqURL[2]) == "":
		forwarding += p.RealPath + "index.html"
	case !p.IsSPA:
		http.Error(w, "dfdsfs", http.StatusNotFound)
		return
	default:
		if reqURL[1] != reqURL[2] {
			forwarding += p.AssetsRealPath + strings.Join(reqURL[2:], "")
			break
		}
		forwarding += p.RealPath + "index.html"
	}
	url, err := url.Parse(forwarding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rr := http.Request{
		Method:        r.Method,
		URL:           url,
		Header:        r.Header,
		Body:          r.Body,
		ContentLength: r.ContentLength,
		Close:         r.Close,
	}
	resp, err := http.DefaultTransport.RoundTrip(&rr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	// TODO: here we assume that there is nothing
	// in our request url after second slash
	// so maybe refactor it somehow
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		http.Error(w, "we need an id in path params", http.StatusBadRequest)
		return
	}
	id := path[len(path)-1]
	if err := h.service.Delete(r.Context(), project.ID(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
