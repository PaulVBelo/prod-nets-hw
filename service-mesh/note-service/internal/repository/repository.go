package repository

import (
	"context"
	"note-service/internal/models"
)

type NoteRepository interface {
  Create(ctx context.Context, note *models.Note) error
  GetByID(ctx context.Context, id string) (*models.Note, error)
  List(ctx context.Context) ([]models.Note, error)
  UpdateDescription(ctx context.Context, id, description string) (*models.Note, error)
  Delete(ctx context.Context, id string) error
}
