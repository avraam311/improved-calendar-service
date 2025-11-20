package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Time("time", time.Now()),
			)
		})
	}
}
