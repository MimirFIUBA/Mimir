package middlewares

import (
	"context"
	"log/slog"
	"net/http"
)

func CreateAPILoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Attach logger to request context
			ctx := context.WithValue(r.Context(), "logger", logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CreateRequestLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			logger.Info(
				"Incoming request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"body", r.Body,
				"status", w.Header().Get("status"),
			)

			// Deletes header after use
			w.Header().Del("status")
		})
	}
}

func ContextWithLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
