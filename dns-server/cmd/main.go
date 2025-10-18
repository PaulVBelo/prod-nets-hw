package main

import (
	"context"
	"dns-server/internal/config"
	dns "dns-server/internal/dns"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	configPath := "config.yaml"

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v", err)
	}

	srv, err := dns.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to init server: %v", err)
	}

	log.Printf("Starting DNS on %s (UDP/TCP). Upstream: %v", cfg.Listen, cfg.Upstream)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("DNS server error: %v", err)
	}
}
