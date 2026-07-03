package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

func InjectLogger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Добавляем logger в контекст запроса
			ctx := context.WithValue(r.Context(), "logger", logger)
			// Передаем дальше с новым контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
