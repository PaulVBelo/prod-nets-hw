package server

import (
	"encoding/xml"
	"net/http"
	"time"

	"note-service/internal/models"
	"note-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	CreateNote            *CreateNoteRequest            `xml:"CreateNoteRequest,omitempty"`
	GetNote               *GetNoteRequest               `xml:"GetNoteRequest,omitempty"`
	ListNotes             *ListNotesRequest             `xml:"ListNotesRequest,omitempty"`
	DeleteNote            *DeleteNoteRequest            `xml:"DeleteNoteRequest,omitempty"`
	UpdateNoteDescription *UpdateNoteDescriptionRequest `xml:"UpdateNoteDescriptionRequest,omitempty"`
}

// === Request types ===

type CreateNoteRequest struct {
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
}

type GetNoteRequest struct {
	ID string `xml:"Id"`
}

type ListNotesRequest struct{}

type DeleteNoteRequest struct {
	ID string `xml:"Id"`
}

type UpdateNoteDescriptionRequest struct {
	ID          string `xml:"Id"`
	Description string `xml:"Description"`
}

// === Response types ===

type CreateNoteResponse struct {
	Note NoteXML `xml:"Note"`
}

type GetNoteResponse struct {
	Note NoteXML `xml:"Note"`
}

type ListNotesResponse struct {
	Notes []NoteXML `xml:"Note"`
}

type DeleteNoteResponse struct {
	OK bool `xml:"Ok"`
}

type UpdateNoteDescriptionResponse struct {
	Note NoteXML `xml:"Note"`
}

// NoteXML — XML-представление заметки
type NoteXML struct {
	ID          string `xml:"Id"`
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	CreatedAt   string `xml:"CreatedAt"`
	UpdatedAt   string `xml:"UpdatedAt"`
}

func noteToXML(n *models.Note) NoteXML {
	return NoteXML{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
		CreatedAt:   n.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   n.UpdatedAt.Format(time.RFC3339),
	}
}

// === Регистрация и обработчик ===

func registerSOAPRoute(r *gin.Engine, noteSvc service.NoteService) {
	r.POST("/soap", soapHandler(noteSvc))
}

func soapHandler(svc service.NoteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env SOAPEnvelope
		decoder := xml.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&env); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid SOAP XML", "details": err.Error()})
			return
		}

		switch {
		case env.Body.CreateNote != nil:
			handleSOAPCreateNote(c, svc, env.Body.CreateNote)
		case env.Body.GetNote != nil:
			handleSOAPGetNote(c, svc, env.Body.GetNote)
		case env.Body.ListNotes != nil:
			handleSOAPListNotes(c, svc, env.Body.ListNotes)
		case env.Body.DeleteNote != nil:
			handleSOAPDeleteNote(c, svc, env.Body.DeleteNote)
		case env.Body.UpdateNoteDescription != nil:
			handleSOAPUpdateDescription(c, svc, env.Body.UpdateNoteDescription)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported SOAP operation"})
		}
	}
}

// === Handlers ===

func handleSOAPCreateNote(c *gin.Context, svc service.NoteService, req *CreateNoteRequest) {
	note := &models.Note{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := svc.Create(c.Request.Context(), note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	respEnv := SOAPEnvelope{
		Body: SOAPBody{
			CreateNote: nil,
		},
	}
	resp := CreateNoteResponse{Note: noteToXML(note)}

	c.XML(http.StatusOK, struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			CreateNoteResponse CreateNoteResponse `xml:"CreateNoteResponse"`
		} `xml:"Body"`
	}{
		Body: struct {
			CreateNoteResponse CreateNoteResponse `xml:"CreateNoteResponse"`
		}{
			CreateNoteResponse: resp,
		},
	})
	_ = respEnv
}

func handleSOAPGetNote(c *gin.Context, svc service.NoteService, req *GetNoteRequest) {
	note, err := svc.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	resp := GetNoteResponse{Note: noteToXML(note)}
	c.XML(http.StatusOK, struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			GetNoteResponse GetNoteResponse `xml:"GetNoteResponse"`
		} `xml:"Body"`
	}{
		Body: struct {
			GetNoteResponse GetNoteResponse `xml:"GetNoteResponse"`
		}{
			GetNoteResponse: resp,
		},
	})
}

func handleSOAPListNotes(c *gin.Context, svc service.NoteService, _ *ListNotesRequest) {
	notes, err := svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notes"})
		return
	}

	resp := ListNotesResponse{
		Notes: make([]NoteXML, 0, len(notes)),
	}
	for i := range notes {
		resp.Notes = append(resp.Notes, noteToXML(&notes[i]))
	}

	c.XML(http.StatusOK, struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			ListNotesResponse ListNotesResponse `xml:"ListNotesResponse"`
		} `xml:"Body"`
	}{
		Body: struct {
			ListNotesResponse ListNotesResponse `xml:"ListNotesResponse"`
		}{
			ListNotesResponse: resp,
		},
	})
}

func handleSOAPDeleteNote(c *gin.Context, svc service.NoteService, req *DeleteNoteRequest) {
	if err := svc.Delete(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	resp := DeleteNoteResponse{OK: true}
	c.XML(http.StatusOK, struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			DeleteNoteResponse DeleteNoteResponse `xml:"DeleteNoteResponse"`
		} `xml:"Body"`
	}{
		Body: struct {
			DeleteNoteResponse DeleteNoteResponse `xml:"DeleteNoteResponse"`
		}{
			DeleteNoteResponse: resp,
		},
	})
}

func handleSOAPUpdateDescription(c *gin.Context, svc service.NoteService, req *UpdateNoteDescriptionRequest) {
	note, err := svc.UpdateDescription(c.Request.Context(), req.ID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note description"})
		return
	}

	resp := UpdateNoteDescriptionResponse{Note: noteToXML(note)}
	c.XML(http.StatusOK, struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			UpdateNoteDescriptionResponse UpdateNoteDescriptionResponse `xml:"UpdateNoteDescriptionResponse"`
		} `xml:"Body"`
	}{
		Body: struct {
			UpdateNoteDescriptionResponse UpdateNoteDescriptionResponse `xml:"UpdateNoteDescriptionResponse"`
		}{
			UpdateNoteDescriptionResponse: resp,
		},
	})
}
