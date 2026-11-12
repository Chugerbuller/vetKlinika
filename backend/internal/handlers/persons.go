package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

func HandlePersons(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r)
	log.Info("Incoming request to HandlePersons")

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
		log.WithField("method", r.Method).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPersonById(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetPersonById")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		log.Warn("Missing person ID in query")
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Warn("Invalid person ID")
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	res, err := st.GetPersonById(id)
	if errors.Is(err, sql.ErrNoRows) {
		log.WithField("id", id).Info("Person not found")
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Error retrieving person from DB")
		http.Error(w, "Error retrieving person", http.StatusInternalServerError)
		return
	}

	log.WithField("id", id).Info("Person retrieved successfully")
	writeJson(w, res)
}

func handleDeletePerson(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleDeletePerson")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		log.Warn("Missing person ID in query")
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Error("Failed to parse person ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// Здесь в вашем коде не было return, оставляю как есть, но id может быть 0
	}

	err = st.DeletePerson(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Error deleting person from DB")
		http.Error(w, "Error deleting person", http.StatusInternalServerError)
		return
	}

	log.WithField("id", id).Info("Person deleted successfully")
	writeJson(w, id)
}

func handleGetPersons(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetPersons")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	res, err := st.GetPersons()
	if err != nil {
		log.WithError(err).Error("Error retrieving persons list from DB")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(*res)).Info("Persons list retrieved successfully")
	writeJson(w, res)
}

func handlePersonByName(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handlePersonByName")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	name := r.URL.Query().Get("name")
	name = strings.ToLower(name)
	if len(name) == 0 {
		log.Warn("Missing name query parameter")
		http.Error(w, "Missing name query parameter", http.StatusBadRequest)
		return
	}

	res, err := st.GetPersonsByName(name)
	if err != nil {
		log.WithError(err).WithField("name", name).Error("Error searching persons by name")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{"name": name, "count": len(*res)}).Info("Search by name completed")
	writeJson(w, res)
}

func handleCreatePerson(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleCreatePerson")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	var person storage.Person
	if err := jsonDeserialize(r.Body, &person); err != nil {
		log.WithError(err).Error("Failed to deserialize person body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := st.AddPerson(person)
	if err != nil {
		log.WithError(err).Error("Error adding person to DB")
		http.Error(w, "Error adding person", http.StatusInternalServerError)
		return
	}

	log.WithField("created_id", id).Info("Person created successfully")
	w.WriteHeader(http.StatusCreated)
	writeJson(w, struct {
		Id int `json:"id"`
	}{
		Id: id,
	})
}

func handleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleUpdatePerson")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	var person storage.Person
	if err := jsonDeserialize(r.Body, &person); err != nil {
		log.WithError(err).Error("Failed to deserialize person body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := st.UpdatePerson(person)
	if err != nil {
		log.WithError(err).WithField("id", person.Id).Error("Error updating person in DB")
		http.Error(w, "Error updating person", http.StatusInternalServerError)
		return
	}

	log.WithField("id", person.Id).Info("Person updated successfully")
	writeJson(w, person.Id)
}
