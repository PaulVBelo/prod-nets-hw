package server

import (
  "mailer-service/internal/config"
  "mailer-service/internal/email"
  "mailer-service/internal/notesclient"
  "net/http"

  "github.com/gin-gonic/gin"
)

type HTTPServer struct {
  cfg         *config.Config
  router      *gin.Engine
  notesClient notesclient.Client
  emailSender email.Sender
}

type sendNoteRequest struct {
  NoteID string `json:"note_id"`
  Email  string `json:"email"`
}

func NewHTTPServer(cfg *config.Config, notes notesclient.Client, email email.Sender) *HTTPServer {
  r := gin.Default()

  s := &HTTPServer{
    cfg:         cfg,
    router:      r,
    notesClient: notes,
    emailSender: email,
  }

  r.POST("/send", s.handleSend)

  return s
}

func (s *HTTPServer) Run() error {
  return s.router.Run(":" + s.cfg.AppPort)
}

func (s *HTTPServer) handleSend(c *gin.Context) {
  var req sendNoteRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
    return
  }

  note, err := s.notesClient.GetNoteByID(req.NoteID)
  if err != nil {
    c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
    return
  }

  if err := s.emailSender.Send(req.Email, note); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email"})
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "status": "sent",
    "to":     req.Email,
    "note":   note,
  })
}
