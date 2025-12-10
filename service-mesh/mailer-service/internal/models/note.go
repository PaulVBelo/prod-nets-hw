package models

type Note struct {
  ID          string `json:"id"`
  Title       string `json:"title"`
  Description string `json:"description"`
  CreatedAt   string `json:"created_at"`
  UpdatedAt   string `json:"updated_at"`
}
