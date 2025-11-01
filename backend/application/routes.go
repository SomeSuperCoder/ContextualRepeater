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
	mux.Handle("/pages/", loadPageRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadPageRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	pageHandler := handlers.PageHandler{
		Repo: repository.NewPageRepo(db),
	}

	mux.HandleFunc("GET /", pageHandler.GetPaged)
	mux.HandleFunc("GET /{id}", pageHandler.Get)
	mux.HandleFunc("POST /", pageHandler.Create)
	mux.HandleFunc("PATCH /{id}", pageHandler.Upadate)
	mux.HandleFunc("DELETE /{id}", pageHandler.Delete)

	return http.StripPrefix("/pages", mux)
}
