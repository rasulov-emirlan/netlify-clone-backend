package project

import (
	"context"
	"io"
)

type Service interface {
	Deploy(ctx context.Context, f io.Reader, name, basePath string, isSPA bool) (Project, error)
	Delete(ctx context.Context, id ID) error
}

type service struct {
}

func NewService() (Service, error) {
	return &service{}, nil
}

func (s *service) Deploy(ctx context.Context, f io.Reader, name, basePath string, isSPA bool) (Project, error) {
	panic("not implemented")
}

func (s *service) Delete(ctx context.Context, id ID) error {
	panic("not implemented")
}
