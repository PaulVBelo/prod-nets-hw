package notesclient

import (
  "crypto/tls"
  "encoding/json"
  "fmt"
  "mailer-service/internal/config"
  "mailer-service/internal/models"
  "net/http"
  "time"
)

type Client interface {
  GetNoteByID(id string) (*models.Note, error)
}

type notesClient struct {
  cfg    *config.Config
  client *http.Client
}

func New(cfg *config.Config) Client {
  transport := &http.Transport{}
  if cfg.SkipTLSVerify {
    transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
  }

  return &notesClient{
    cfg: cfg,
    client: &http.Client{
      Timeout:   2 * time.Second,
      Transport: transport,
    },
  }
}

func (c *notesClient) GetNoteByID(id string) (*models.Note, error) {
  url := fmt.Sprintf("%s:%s/notes/%s", c.cfg.NotesBaseURL, c.cfg.NotesBasePort, id)

  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("cannot build request: %w", err)
  }

  resp, err := c.client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("error calling notes-service: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("notes-service returned status %d", resp.StatusCode)
  }

  var note models.Note
  if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
    return nil, fmt.Errorf("decode error: %w", err)
  }

  return &note, nil
}
