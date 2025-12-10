package main

import (
  "log"

  "mailer-service/internal/config"
  "mailer-service/internal/email"
  "mailer-service/internal/notesclient"
  "mailer-service/internal/server"
)

func main() {
  cfg := config.Load()

  notesClient := notesclient.New(cfg)

  emailSender := email.NewSender(cfg)

  httpServer := server.NewHTTPServer(cfg, notesClient, emailSender)

  log.Printf("starting mailer-service on :%s", cfg.AppPort)

  if err := httpServer.Run(); err != nil {
    log.Fatalf("server stopped: %v", err)
  }
}
