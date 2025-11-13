package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
)

type SentenceHandler struct {
	Repo *repository.SentenceRepo
}

var sentenceArrayPath = []string{"sentences"}

func (h *SentenceHandler) Push(w http.ResponseWriter, r *http.Request) {
	Push(w, r, h.Repo, func(r *models.SentenceCreateRequest) *models.Sentence {
		return &models.Sentence{
			MainContent:  r.MainContent,
			ExtraContent: r.ExtraContent,
			Reviews:      make([]models.Review, 0),
		}
	}, sentenceArrayPath)
}

func (h *SentenceHandler) Pull(w http.ResponseWriter, r *http.Request) {
	Pull(w, r, h.Repo, sentenceArrayPath)
}

func (h *SentenceHandler) ArrayUpdate(w http.ResponseWriter, r *http.Request) {
	ArrayUpdate(w, r, h.Repo, sentenceArrayPath)
}
