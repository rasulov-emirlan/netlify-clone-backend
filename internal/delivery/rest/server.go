package rest

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/gommon/color"
)

type server struct {
	listener *http.Server
	domain   string

	projectHandler http.Handler
}

func NewServer(
	port, domain string,
	timeR, timeW time.Duration,
	pHandler http.Handler,
) (*server, error) {
	if pHandler == nil {
		return nil, errors.New("server: handlers can't be nil")
	}
	s := &http.Server{
		ReadTimeout:  timeR,
		WriteTimeout: timeW,
		Addr:         port,
	}
	return &server{
		listener:       s,
		domain:         domain,
		projectHandler: pHandler,
	}, nil
}

func (s *server) Start() error {
	h, err := s.registerRoutes()
	if err != nil {
		return err
	}
	s.listener.Handler = h
	log.Println(color.Green("Server is listening at port:" + s.listener.Addr))
	return s.listener.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.listener.Shutdown(ctx)
}

func (s *server) registerRoutes() (http.Handler, error) {
	return s.projectHandler, nil
}
