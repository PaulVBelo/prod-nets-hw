package repository

import (
	"context"
	"fmt"

	"note-service/internal/config"
	"note-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormNoteRepository struct {
	db *gorm.DB
}

func NewGormNoteRepository(db *gorm.DB) NoteRepository {
	return &GormNoteRepository{db: db}
}

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Note{}); err != nil {
		return nil, err
	}

	return db, nil
}

func NewNoteRepositoryFromConfig(cfg *config.Config) (NoteRepository, error) {
	shardHosts := cfg.DBShardHosts()
	if len(shardHosts) == 0 {
		db, err := NewGormDB(cfg)
		if err != nil {
			return nil, err
		}
		return NewGormNoteRepository(db), nil
	}

	// shard mode
	var shards []NoteRepository
	for _, host := range shardHosts {
		localCfg := *cfg
		localCfg.DBHost = host
		db, err := NewGormDB(&localCfg)
		if err != nil {
			return nil, err
		}
		shards = append(shards, NewGormNoteRepository(db))
	}

	return NewShardedNoteRepository(shards), nil
}

func (r *GormNoteRepository) Create(ctx context.Context, note *models.Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *GormNoteRepository) GetByID(ctx context.Context, id string) (*models.Note, error) {
	var note models.Note
	if err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormNoteRepository) List(ctx context.Context) ([]models.Note, error) {
	var notes []models.Note
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormNoteRepository) UpdateDescription(ctx context.Context, id, description string) (*models.Note, error) {
	var note models.Note
	if err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	note.Description = description
	if err := r.db.WithContext(ctx).Save(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormNoteRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Note{}, "id = ?", id).Error
}
