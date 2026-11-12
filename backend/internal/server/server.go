package server

import (
	"vetklinika/backend/internal/handlers"
	"vetklinika/backend/internal/middleware"
	"vetklinika/backend/internal/storage"

	"net"
	"net/http"

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
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	if ready != nil {
		close(ready)
	}
	err = http.Serve(l, s.router)
	if err == http.ErrServerClosed {
		s.logger.WithError(err).Info("server gracefully closed")
	} else if err != nil {
		s.logger.WithError(err).Fatal("error running server")
	}
	return nil
}
func serveStaticDir(mux *http.ServeMux, pathPrefix, dirPath string) {
	fs := http.StripPrefix(
		pathPrefix,
		http.FileServer(http.Dir(dirPath)),
	)
	mux.Handle(pathPrefix, fs)
}
