package server

import (
	"net/http"
	"vetklinika/backend/internal/handlers"
	"vetklinika/backend/internal/middleware"
	"vetklinika/backend/internal/storage"
)

type Server struct {
	router  *http.ServeMux
	storage *storage.Storage
}

func NewServer(st *storage.Storage) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		storage: st,
	}
	s.initEndpoints()
	return s
}
func (s *Server) initEndpoints() {
	mw := middleware.InjectStorage(s.storage)
	s.router.Handle("/person/create/", http.StripPrefix("/person/create/", http.FileServer(http.Dir("./frontend/create-person"))))
	s.router.Handle("/person/", http.StripPrefix("/person/", http.FileServer(http.Dir("./frontend/person-card"))))
	s.router.Handle("/", http.FileServer(http.Dir("./frontend/src")))
	s.router.Handle("/api/pets", mw(http.HandlerFunc(handlers.HandlePets)))
	s.router.Handle("/api/persons", mw(http.HandlerFunc(handlers.HandlePersons)))
	s.router.Handle("/api/types", mw(http.HandlerFunc(handlers.HandleTypes)))
}
func (s *Server) Start(port string) error {
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		return err
	}
	return nil
}
