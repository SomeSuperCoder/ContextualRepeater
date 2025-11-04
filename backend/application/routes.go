package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/handlers"
	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func loadRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// =========================
	// Page routes
	// =========================
	pageHandler := handlers.PageHandler{
		Repo: repository.NewPageRepo(db),
	}

	mux.HandleFunc("GET /pages/{$}", pageHandler.GetPaged)
	mux.HandleFunc("GET /pages/{id}", pageHandler.Get)
	mux.HandleFunc("POST /pages/{$}", pageHandler.Create)
	mux.HandleFunc("PATCH /pages/{id}", pageHandler.Upadate)
	mux.HandleFunc("DELETE /pages/{id}", pageHandler.Delete)

	// =========================
	// Sentence routes
	// =========================
	sentenceHandler := handlers.SentenceHandler{
		Repo: repository.NewSentenceRepo(db),
	}

	mux.HandleFunc("POST /pages/{id}/sentences/", sentenceHandler.Push)
	mux.HandleFunc("PATCH /pages/{id}/sentences/", sentenceHandler.ArrayUpdate)
	mux.HandleFunc("DELETE /pages/{id}/sentences/", sentenceHandler.Pull)

	return middleware.LoggerMiddleware(mux)
}
