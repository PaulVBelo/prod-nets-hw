package repository

import (
	"context"
	"hash/fnv"

	"note-service/internal/models"

	"github.com/google/uuid"
)

type shardedNoteRepository struct {
	shards []NoteRepository
}

func NewShardedNoteRepository(shards []NoteRepository) NoteRepository {
	return &shardedNoteRepository{shards: shards}
}

func (r *shardedNoteRepository) shardIndexByID(id string) int {
	if len(r.shards) == 1 {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))
	return int(h.Sum32() % uint32(len(r.shards)))
}

func (r *shardedNoteRepository) shardByNote(note *models.Note) NoteRepository {
	if note.ID == "" {
		return r.shards[0]
	}
	idx := r.shardIndexByID(note.ID)
	return r.shards[idx]
}

func (r *shardedNoteRepository) shardByID(id string) NoteRepository {
	idx := r.shardIndexByID(id)
	return r.shards[idx]
}

func (r *shardedNoteRepository) Create(ctx context.Context, note *models.Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}
	shard := r.shardByNote(note)
	return shard.Create(ctx, note)
}

func (r *shardedNoteRepository) GetByID(ctx context.Context, id string) (*models.Note, error) {
	shard := r.shardByID(id)
	return shard.GetByID(ctx, id)
}

func (r *shardedNoteRepository) List(ctx context.Context) ([]models.Note, error) {
	var all []models.Note
	for _, shard := range r.shards {
		notes, err := shard.List(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, notes...)
	}
	return all, nil
}

func (r *shardedNoteRepository) UpdateDescription(ctx context.Context, id, description string) (*models.Note, error) {
	shard := r.shardByID(id)
	return shard.UpdateDescription(ctx, id, description)
}

func (r *shardedNoteRepository) Delete(ctx context.Context, id string) error {
	shard := r.shardByID(id)
	return shard.Delete(ctx, id)
}
