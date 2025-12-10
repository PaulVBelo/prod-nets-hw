package main

import (
	"log"

	"note-service/internal/config"
	"note-service/internal/repository"
	"note-service/internal/server"
	"note-service/internal/service"

	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.MustLoad(config.NewEnvLoader())

	noteRepo, err := repository.NewNoteRepositoryFromConfig(cfg)
	if err != nil {
		log.Fatalf("failed tp init repository: %v", err)
	}

	noteSvc := service.NewNoteService(noteRepo)

	httpSrv := server.NewGinServer(cfg, noteSvc)
	grpcSrv := server.NewGRPCServer(cfg, noteSvc)

	var g errgroup.Group

	g.Go(func() error {
		return httpSrv.Run()
	})

	g.Go(func() error {
		return grpcSrv.Run()
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
