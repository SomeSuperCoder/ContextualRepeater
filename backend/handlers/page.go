package handlers

import (
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PageHandler struct {
	Repo *repository.PageRepo
}

func (h *PageHandler) GetPaged(r *http.Request, w http.ResponseWriter) {
	FindPaged(w, r, h.Repo)
}

func (h *PageHandler) Get(r *http.Request, w http.ResponseWriter) {
	GetByID(w, r, h.Repo)
}

func (h *PageHandler) Create(r *http.Request, w http.ResponseWriter) {
	Create(w, r, h.Repo, func(r *models.PageCreateRequest) *models.Page {
		return &models.Page{
			Title:       r.Title,
			Language:    r.Language,
			Sentences:   make([]models.Sentence, 0),
			ReviewsDone: 0,
			ReviewDate:  time.Now(),
		}
	}, func(createdID bson.ObjectID) {})
}

func (h *PageHandler) Upadate(r *http.Request, w http.ResponseWriter) {
	DefaultUpdate(w, r, h.Repo)
}

func (h *PageHandler) Delete(r *http.Request, w http.ResponseWriter) {
	Delete(w, r, h.Repo)
}
