package email

import (
  "fmt"
  "log"
  "mailer-service/internal/config"
  "mailer-service/internal/models"
  "net/smtp"
)

type Sender interface {
  Send(to string, note *models.Note) error
}

type smtpSender struct {
  cfg *config.Config
}

func NewSender(cfg *config.Config) Sender {
  return &smtpSender{cfg: cfg}
}

func (s *smtpSender) Send(to string, note *models.Note) error {
  addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

  var auth smtp.Auth
  if s.cfg.SMTPUser != "" {
    auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
  }

  subject := fmt.Sprintf("Note %s: %s", note.ID, note.Title)
  body := fmt.Sprintf(
    "Note ID: %s\nTitle: %s\nDescription:\n%s\nCreated: %s\nUpdated: %s\n",
    note.ID,
    note.Title,
    note.Description,
    note.CreatedAt,
    note.UpdatedAt,
  )

  msg := []byte(fmt.Sprintf(
    "From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
    s.cfg.SMTPFrom, to, subject, body,
  ))

  log.Printf("Sending email → %s (via %s)", to, addr)

  return smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, msg)
}
