package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		fmt.Printf("%s %s took %v - %d %s\n", r.Method, r.URL.Path, duration, wrapped.statusCode, http.StatusText(wrapped.statusCode))

		if wrapped.statusCode/100 != 2 {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal("Failed to read response body in the logger middleware")
			}
			fmt.Print(string(data))
		}
	})
}

// Interceptor sturct
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Intercept the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
