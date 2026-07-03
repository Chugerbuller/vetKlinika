package server

import (
	"net"
	"net/http"
	"strings"
	"vetklinika/backend/internal/handlers"
	"vetklinika/backend/internal/middleware"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

type Server struct {
	router  *http.ServeMux
	storage *storage.Storage
	logger  *logrus.Logger
}

func NewServer(st *storage.Storage, logs *logrus.Logger) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		storage: st,
		logger:  logs,
	}
	s.initEndpoints()
	return s
}
func (s *Server) initEndpoints() {
	mw := middleware.NewMw(s.storage, s.logger)
	//handlers.HandlePets
	s.router.Handle("/api/pets", middleware.Chain(handlers.HandlePets, mw...))
	s.router.Handle("/api/persons", middleware.Chain(handlers.HandlePersons, mw...))
	s.router.Handle("/api/types", middleware.Chain(handlers.HandleTypes, mw...))
	s.router.Handle("/api/appointments", middleware.Chain(handlers.HandleAppointments, mw...))
	s.router.Handle("/api/documents", middleware.Chain(handlers.HandleDocuments, mw...))

	serveStaticDir(s.router, "/person/create/", "./frontend/create-person")
	serveStaticDir(s.router, "/person/", "./frontend/person-card")
	serveStaticDir(s.router, "/pets/", "./frontend/pet-card")
	serveStaticDir(s.router, "/visit/create/", "./frontend/create-appointment")
	serveStaticDir(s.router, "/visit/", "./frontend/appointment-card")

	s.router.Handle("/", http.FileServer(http.Dir("./frontend/src")))

}
func (s *Server) Start(port string, ready chan<- struct{}) error {
	// 1. Пытаемся открыть порт
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// 2. Порт успешно открыт!
	// Закрываем канал, чтобы уведомить тех, кто ждет
	if ready != nil {
		close(ready)
	}

	// 3. Запускаем сервер на этом слушателе
	// Это блокирующая операция, как и ListenAndServe
	return http.Serve(l, s.router)
}
func serveStaticDir(mux *http.ServeMux, pathPrefix string, dirPath string) {
	fs := http.FileServer(http.Dir(dirPath))
	// Добавляем слеш в конец, если его нет, для корректной работы StripPrefix
	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}
	mux.Handle(pathPrefix, http.StripPrefix(pathPrefix, fs))
}
