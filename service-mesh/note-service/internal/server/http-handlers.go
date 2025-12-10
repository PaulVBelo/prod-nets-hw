package server

import (
  "net/http"

  "note-service/internal/models"
  "note-service/internal/service"

  "github.com/gin-gonic/gin"
)

func registerNoteRoutes(r *gin.Engine, noteSvc service.NoteService) {
  notes := r.Group("/notes")
  {
    notes.POST("", createNoteHandler(noteSvc))
    notes.GET("", listNotesHandler(noteSvc))
    notes.GET("/:id", getNoteHandler(noteSvc))
    notes.PATCH("/:id/description", updateNoteDescriptionHandler(noteSvc))
    notes.DELETE("/:id", deleteNoteHandler(noteSvc))
  }
}

type createNoteRequest struct {
  Title       string `json:"title" binding:"required"`
  Description string `json:"description"`
}

func createNoteHandler(svc service.NoteService) gin.HandlerFunc {
  return func(c *gin.Context) {
    var req createNoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    note := &models.Note{
      Title:       req.Title,
      Description: req.Description,
    }

    if err := svc.Create(c.Request.Context(), note); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
      return
    }

    c.JSON(http.StatusCreated, note)
  }
}

func listNotesHandler(svc service.NoteService) gin.HandlerFunc {
  return func(c *gin.Context) {
    notes, err := svc.List(c.Request.Context())
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notes"})
      return
    }
    c.JSON(http.StatusOK, notes)
  }
}

func getNoteHandler(svc service.NoteService) gin.HandlerFunc {
  return func(c *gin.Context) {
    id := c.Param("id")
    note, err := svc.GetByID(c.Request.Context(), id)
    if err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
      return
    }
    c.JSON(http.StatusOK, note)
  }
}

type updateDescriptionRequest struct {
  Description string `json:"description" binding:"required"`
}

func updateNoteDescriptionHandler(svc service.NoteService) gin.HandlerFunc {
  return func(c *gin.Context) {
    id := c.Param("id")
    var req updateDescriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    note, err := svc.UpdateDescription(c.Request.Context(), id, req.Description)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
      return
    }
    c.JSON(http.StatusOK, note)
  }
}

func deleteNoteHandler(svc service.NoteService) gin.HandlerFunc {
  return func(c *gin.Context) {
    id := c.Param("id")
    if err := svc.Delete(c.Request.Context(), id); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
      return
    }
    c.Status(http.StatusNoContent)
  }
}
