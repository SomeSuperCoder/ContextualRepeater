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

func (h *PageHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
	FindPaged(w, r, h.Repo)
}

func (h *PageHandler) Get(w http.ResponseWriter, r *http.Request) {
	GetByID(w, r, h.Repo)
}

func (h *PageHandler) Create(w http.ResponseWriter, r *http.Request) {
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

func (h *PageHandler) Upadate(w http.ResponseWriter, r *http.Request) {
	DefaultUpdate(w, r, h.Repo)
}

func (h *PageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, h.Repo)
}
