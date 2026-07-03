package middleware

import (
	"net/http"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

func Chain(h http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
	// Теперь h — это функция, мы можем обернуть её прямо здесь один раз
	var finalHandler http.Handler = h

	// Применяем middleware в обратном порядке,
	// чтобы первый middleware в списке выполнялся первым
	for i := len(mws) - 1; i >= 0; i-- {
		finalHandler = mws[i](finalHandler)
	}

	return finalHandler
}
func NewMw(st *storage.Storage, logger *logrus.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		InjectLogger(logger),
		InjectStorage(st),
	}
}
