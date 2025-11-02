package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
)

type SentenceHandler struct {
	Repo *repository.SentenceRepo
}

func (h *SentenceHandler) Push(w http.ResponseWriter, r *http.Request) {
	Push(w, r, h.Repo, func(r *models.SentenceCreateRequest) *models.Sentence {
		return &models.Sentence{
			MainContent:  r.MainContent,
			ExtraContent: r.ExtraContent,
			Reviews:      make([]models.Review, 0),
		}
	})
}
