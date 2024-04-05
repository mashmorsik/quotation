package middleware

import (
	"github.com/mashmorsik/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Infof("Started %s %s", r.Method, r.URL.Path)

		handler.ServeHTTP(w, r)

		logger.Infof("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}
