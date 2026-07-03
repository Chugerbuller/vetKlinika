package handlers

import (
	"net/http"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

func getLogger(r *http.Request) *logrus.Entry {
	l, ok := r.Context().Value("logger").(*logrus.Logger)
	if !ok {
		// Если логгера нет в контексте, возвращаем стандартный, чтобы не паниковать
		return logrus.NewEntry(logrus.StandardLogger())
	}
	// Возвращаем Entry с предустановленными полями запроса
	return l.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"remote": r.RemoteAddr,
	})
}

// Вспомогательная функция для извлечения хранилища
func getStorage(r *http.Request) *storage.Storage {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		return nil
	}
	return st
}
