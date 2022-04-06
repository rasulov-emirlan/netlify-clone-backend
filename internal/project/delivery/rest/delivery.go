package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
)

type Handler struct {
	service project.Service
}

func NewHandler(s project.Service) (*Handler, error) {
	if s == nil {
		return nil, errors.New("project: arguments for NewHandler can't be nil")
	}
	return &Handler{
		service: s,
	}, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
	fmt.Println(req.RequestURI)
	fmt.Println(req.Proto)
	fmt.Println(req.Proto)
}
