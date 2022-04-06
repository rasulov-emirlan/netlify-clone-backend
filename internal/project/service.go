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

type Logger interface {
	Write(errLevel int, format string, v ...interface{}) error
}

const (
	Info = iota
	Error
	Trace
	Panic
)

type FileSystem interface {
	Upload(ctx context.Context, f io.Reader, id string) (path string, err error)
	Delete(ctx context.Context, id string) error
}

type Repository interface {
	Create(ctx context.Context, p Project) (Project, error)
	Read(ctx context.Context, id ID) (Project, error)

	Update(ctx context.Context, id ID, p Project) (Project, error)
	Delete(ctx context.Context, id ID) error
}

type service struct {
	fs   FileSystem
	repo Repository
	log  Logger
}

func NewService(fs FileSystem, repo Repository, log Logger) (Service, error) {
	if fs == nil || repo == nil || log == nil {
		return nil, errors.New("project: arguments for creating new service can't be nil")
	}
	return &service{
		fs:   fs,
		repo: repo,
		log:  log,
	}, nil
}

func (s *service) Deploy(ctx context.Context, f io.Reader, name, basePath string, isSPA bool) (Project, error) {
	p, err := NewModel(name, basePath, isSPA)
	if err != nil {
		return p, err
	}
	p, err = s.repo.Create(ctx, p)
	if err != nil {
		return p, err
	}
	path, err := s.fs.Upload(ctx, f, name)
	if err != nil {
		return p, err
	}
	p.RealPath = path
	return p, nil
}

func (s *service) Delete(ctx context.Context, id ID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return s.fs.Delete(ctx, string(id))
}
