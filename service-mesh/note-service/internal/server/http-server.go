package server

import (
  "fmt"
  "log"
  "net/http"

  "note-service/internal/config"
  "note-service/internal/service"

  "github.com/gin-gonic/gin"
)

type HTTPServer interface {
  Run() error
}

type GinServer struct {
  engine *gin.Engine
  cfg    *config.Config
}

func NewGinServer(cfg *config.Config, noteSvc service.NoteService) HTTPServer {
  if cfg.GinMode != "" {
    gin.SetMode(cfg.GinMode)
  }

  router := gin.New()
  router.Use(gin.Logger())
  router.Use(gin.Recovery())

  router.GET("/healthz", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
  })

  registerNoteRoutes(router, noteSvc)
  registerSOAPRoute(router, noteSvc)

  return &GinServer{
    engine: router,
    cfg:    cfg,
  }
}

func (s *GinServer) Run() error {
  addr := fmt.Sprintf(":%s", s.cfg.AppPort)
  log.Printf("starting HTTP server on %s", addr)
  return s.engine.Run(addr)
}
