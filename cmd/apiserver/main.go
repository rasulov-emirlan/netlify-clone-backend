package main

import (
	"log"
	"os"

	"github.com/rasulov-emirlan/netlify-clone-backend/config"
	"github.com/rasulov-emirlan/netlify-clone-backend/internal/delivery/rest"
	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
	projectR "github.com/rasulov-emirlan/netlify-clone-backend/internal/project/delivery/rest"
	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project/fs/miniofs"
	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project/store/postgres"
	"github.com/rasulov-emirlan/netlify-clone-backend/pkg/db"
)

var configfiles []string

func main() {
	if len(os.Args) > 1 {
		configfiles = append(configfiles, os.Args[1:]...)
	}
	cfg, err := config.NewConfig(configfiles...)
	if err != nil {
		log.Fatal(err)
	}
	dbConn, err := db.NewGORM(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	filesystem, err := miniofs.NewFileSystem(cfg.Minio.URL, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, cfg.Minio.UseSSL)
	if err != nil {
		log.Fatal(err)
	}
	repo, err := postgres.NewRepo(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	service, err := project.NewService(filesystem, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	if service == nil {
		log.Fatal(service)
	}
	router, err := projectR.NewHandler(service)
	if err != nil {
		log.Fatal(err)
	}
	s, err := rest.NewServer(cfg.Server.Port, cfg.Server.Domain, cfg.Server.TimeoutRead, cfg.Server.TimeoutWrite, router)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start())
}
