package project

import (
	"context"
	"errors"
	"io"
)

type Service interface {
	Deploy(ctx context.Context, f io.Reader, name, basePath string, isSPA bool) (Project, error)
	Delete(ctx context.Context, id ID) error
}

type FileSystem interface {
	Upload(ctx context.Context, f io.Reader, id string) (path string, err error)
	Delete(ctx context.Context, id string) error
}

type Repository interface {
	Create(ctx context.Context, p Project) (Project, error)
	Read(ctx context.Context, id ID) (Project, error)

	Update(ctx context.Context, id ID, p Project) (Project, error)
	Delete(ctx context.Context, p Project) error
}

type service struct {
	fs   FileSystem
	repo Repository
}

func NewService(fs FileSystem, repo Repository) (Service, error) {
	if fs == nil || repo == nil {
		return nil, errors.New("project: arguments for creating new service can't be nil")
	}
	return &service{
		fs:   fs,
		repo: repo,
	}, nil
}

func (s *service) Deploy(ctx context.Context, f io.Reader, name, basePath string, isSPA bool) (Project, error) {
	panic("not implemented")
}

func (s *service) Delete(ctx context.Context, id ID) error {
	panic("not implemented")
}
