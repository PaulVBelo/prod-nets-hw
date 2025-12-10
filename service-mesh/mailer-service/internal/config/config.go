package config

import (
  "os"
  "strconv"
)

type Config struct {
  AppPort string

  NotesBaseURL  string
  NotesBasePort string
  SkipTLSVerify bool

  SMTPHost string
  SMTPPort int
  SMTPUser string
  SMTPPass string
  SMTPFrom string
}

func Load() *Config {
  return &Config{
    AppPort:       getenv("APP_PORT", "8080"),
    NotesBaseURL:  getenv("NOTES_BASE_URL", "https://notes-sidecar"),
    NotesBasePort: getenv("NOTES_BASE_PORT", "443"),
    SkipTLSVerify: getenv("SKIP_TLS_VERIFY", "false") == "true",

    SMTPHost: getenv("SMTP_HOST", "mailhog"),
    SMTPPort: getenvInt("SMTP_PORT", 1025),
    SMTPUser: getenv("SMTP_USER", ""),
    SMTPPass: getenv("SMTP_PASS", ""),
    SMTPFrom: getenv("SMTP_FROM", "notes@example.com"),
  }
}

func getenv(key, def string) string {
  if v := os.Getenv(key); v != "" {
    return v
  }
  return def
}

func getenvInt(key string, def int) int {
  if v := os.Getenv(key); v != "" {
    if n, err := strconv.Atoi(v); err == nil {
      return n
    }
  }
  return def
}
