package server

import (
  "context"
  "time"

  "note-service/api/proto"
  "note-service/internal/models"
  "note-service/internal/service"
)

type noteHandler struct {
  gen.UnimplementedNotesServiceServer
  svc service.NoteService
}

func (h *noteHandler) CreateNote(ctx context.Context, req *gen.CreateNoteRequest) (*gen.NoteResponse, error) {
  note := &models.Note{
    Title:       req.GetTitle(),
    Description: req.GetDescription(),
  }
  if err := h.svc.Create(ctx, note); err != nil {
    return nil, err
  }
  return &gen.NoteResponse{Note: toProto(note)}, nil
}

func (h *noteHandler) ListNotes(ctx context.Context, _ *gen.ListNotesRequest) (*gen.ListNotesResponse, error) {
  notes, err := h.svc.List(ctx)
  if err != nil {
    return nil, err
  }

  res := &gen.ListNotesResponse{
    Notes: make([]*gen.Note, 0, len(notes)),
  }
  for i := range notes {
    res.Notes = append(res.Notes, toProto(&notes[i]))
  }
  return res, nil
}

func (h *noteHandler) GetNote(ctx context.Context, req *gen.GetNoteRequest) (*gen.NoteResponse, error) {
  note, err := h.svc.GetByID(ctx, req.GetId())
  if err != nil {
    return nil, err
  }
  return &gen.NoteResponse{Note: toProto(note)}, nil
}

func (h *noteHandler) DeleteNote(ctx context.Context, req *gen.DeleteNoteRequest) (*gen.DeleteNoteResponse, error) {
  if err := h.svc.Delete(ctx, req.GetId()); err != nil {
    return nil, err
  }
  return &gen.DeleteNoteResponse{Ok: true}, nil
}

func (h *noteHandler) UpdateNoteDescription(ctx context.Context, req *gen.UpdateNoteDescriptionRequest) (*gen.NoteResponse, error) {
  note, err := h.svc.UpdateDescription(ctx, req.GetId(), req.GetDescription())
  if err != nil {
    return nil, err
  }
  return &gen.NoteResponse{Note: toProto(note)}, nil
}

func toProto(note *models.Note) *gen.Note {
  return &gen.Note{
    Id:          note.ID,
    Title:       note.Title,
    Description: note.Description,
    CreatedAt:   note.CreatedAt.Format(time.RFC3339),
    UpdatedAt:   note.UpdatedAt.Format(time.RFC3339),
  }
}
