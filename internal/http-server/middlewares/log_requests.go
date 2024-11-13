package middlewares

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func LogRequestsInfo(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("request", "method", r.Method, "url", r.URL.String())
			next.ServeHTTP(w, r)
		})
	}
}
