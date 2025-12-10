package server

import (
  "fmt"
  "log"
  "net"

  "note-service/api/proto"
  "note-service/internal/config"
  "note-service/internal/service"

  "google.golang.org/grpc"
)

type GRPCServer interface {
  Run() error
}

type grpcServer struct {
  cfg     *config.Config
  svc     service.NoteService
  server  *grpc.Server
}

func NewGRPCServer(cfg *config.Config, noteSvc service.NoteService) GRPCServer {
  s := grpc.NewServer()
  gen.RegisterNotesServiceServer(s, &noteHandler{
    svc: noteSvc,
  })

  return &grpcServer{
    cfg:    cfg,
    svc:    noteSvc,
    server: s,
  }
}

func (g *grpcServer) Run() error {
  addr := fmt.Sprintf(":%s", g.cfg.GRPCPort)
  lis, err := net.Listen("tcp", addr)
  if err != nil {
    return fmt.Errorf("failed to listen on %s: %w", addr, err)
  }
  log.Printf("starting gRPC server on %s", addr)
  return g.server.Serve(lis)
}
