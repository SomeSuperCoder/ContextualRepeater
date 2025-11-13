package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/repository"
)

type ReviewHandler struct {
	Repo *repository.ReviewRepo
}

var reviewArrayPath = []string{"sentences", "reviews"}

func (h *ReviewHandler) Push(w http.ResponseWriter, r *http.Request) {
	Push(w, r, h.Repo, func(r *string) *string {
		return r
	}, reviewArrayPath)
}

func (h *ReviewHandler) Pull(w http.ResponseWriter, r *http.Request) {
	Pull(w, r, h.Repo, reviewArrayPath)
}

func (h *ReviewHandler) ArrayUpdate(w http.ResponseWriter, r *http.Request) {
	ArrayUpdate(w, r, h.Repo, reviewArrayPath)
}
