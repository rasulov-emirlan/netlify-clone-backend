package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	projectRouter "github.com/rasulov-emirlan/netlify-clone-backend/internal/project/delivery/rest"
)

type server struct {
	listener *http.Server
	domain   string

	projectHandler *projectRouter.Handler
}

func NewServer(
	port, domain string,
	timeR, timeW time.Duration,
	pHandler *projectRouter.Handler,
) (*server, error) {
	s := &http.Server{
		ReadTimeout:  timeR,
		WriteTimeout: timeW,
		Addr:         port,
	}
	return &server{
		listener: s,
		domain:   domain,
	}, nil
}

func (s *server) Start() error {
	h, err := s.registerRoutes()
	if err != nil {
		return err
	}
	s.listener.Handler = h
	return s.listener.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.listener.Shutdown(ctx)
}

func (s *server) registerRoutes() (*mux.Router, error) {
	m := mux.NewRouter()
	m.Handle("/", s.projectHandler)
	return m, nil
}
