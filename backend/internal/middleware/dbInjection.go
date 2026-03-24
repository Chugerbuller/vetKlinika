package middleware

import (
	"context"
	"net/http"
	"vetklinika/backend/internal/storage"
)

func InjectStorage(st *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Добавляем db в контекст запроса
			ctx := context.WithValue(r.Context(), "storage", st)
			// Передаем дальше с новым контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
