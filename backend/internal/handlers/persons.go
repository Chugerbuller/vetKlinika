package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"vetklinika/backend/internal/storage"
)

func HandlePersons(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("name") {
			handlePersonByName(w, r)
		} else if r.URL.Query().Has("id") {
			handleGetPersonById(w, r)
		} else {
			handleGetPersons(w, r)
		}
	case http.MethodPost:
		handleCreatePerson(w, r)
	case http.MethodPut:
		handleUpdatePerson(w, r)
	case http.MethodDelete:
		handleDeletePerson(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPersonById(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := st.GetPersonById(id)
	if err != nil {
		http.Error(w, "Error retrieving person", http.StatusInternalServerError)
		return
	}
	writeJson(w, res)
}

func handleDeletePerson(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = st.DeletePerson(id)
	if err != nil {
		http.Error(w, "Error deleting person", http.StatusInternalServerError)
		return
	}
	writeJson(w, id)
}
func handleGetPersons(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	res, err := st.GetPersons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJson(w, res)
}
func handlePersonByName(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	name := r.URL.Query().Get("name")
	name = strings.ToLower(name)
	if len(name) == 0 {
		http.Error(w, "Missing name query parameter", http.StatusBadRequest)
		return
	}
	res, err := st.GetPersonsByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJson(w, res)
}
func handleCreatePerson(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	var person storage.Person
	if err := jsonDeserialize(r.Body, &person); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := st.AddPerson(person)
	if err != nil {
		http.Error(w, "Error adding person", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJson(w, struct {
		Id int `json:"id"`
	}{
		Id: id,
	})

}
func handleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	var person storage.Person
	if err := jsonDeserialize(r.Body, &person); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := st.UpdatePerson(person)
	if err != nil {
		http.Error(w, "Error updating person", http.StatusInternalServerError)
		return
	}
	writeJson(w, person.Id)
}
