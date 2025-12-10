package service

import (
	"context"

	"note-service/internal/models"
	"note-service/internal/repository"
)

type NoteService interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id string) (*models.Note, error)
	List(ctx context.Context) ([]models.Note, error)
	UpdateDescription(ctx context.Context, id, description string) (*models.Note, error)
	Delete(ctx context.Context, id string) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) Create(ctx context.Context, note *models.Note) error {
	return s.repo.Create(ctx, note)
}

func (s *noteService) GetByID(ctx context.Context, id string) (*models.Note, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *noteService) List(ctx context.Context) ([]models.Note, error) {
	return s.repo.List(ctx)
}

func (s *noteService) UpdateDescription(ctx context.Context, id, description string) (*models.Note, error) {
	return s.repo.UpdateDescription(ctx, id, description)
}

func (s *noteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
