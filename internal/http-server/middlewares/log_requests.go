package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func LogRequestsInfo(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := time.Now()
			next.ServeHTTP(w, r)
			log.Info("request", "method", r.Method, "url", r.URL.String(), "time(nanoseconds)", time.Since(current).Microseconds())
		})
	}
}
