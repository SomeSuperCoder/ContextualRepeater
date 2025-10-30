package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func loadRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	return middleware.LoggerMiddleware(mux)
}
