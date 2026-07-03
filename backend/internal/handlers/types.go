package handlers

import (
	"net/http"
	"strconv"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

func HandleTypes(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r)
	log.Info("Incoming request to HandleTypes")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			HandleGetTypeById(w, r)
		} else if r.URL.Query().Has("name") {
			HandleGetTypeByName(w, r)
		} else {
			HandleGetTypes(w, r)
		}
	case http.MethodPost:
		HandleCreateType(w, r)
	default:
		log.WithField("method", r.Method).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleGetTypes(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleGetTypes")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	types, err := st.GetTypes()
	if err != nil {
		log.WithError(err).Error("Failed to get types from DB")
		http.Error(w, "Failed to get types", http.StatusInternalServerError)
		return
	}

	if types == nil {
		log.Info("No types found, returning empty struct")
		writeJson(w, struct{}{})
	} else {
		log.WithField("count", len(*types)).Info("Types retrieved successfully")
		writeJson(w, types)
	}
}

func HandleGetTypeById(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleGetTypeById")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		log.Warn("Missing type ID in query")
		http.Error(w, "Missing type ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Warn("Invalid type ID format")
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	log = log.WithField("type_id", id)
	t, err := st.GetTypeById(id)
	if err != nil {
		log.WithError(err).Error("Failed to get type by ID from DB")
		http.Error(w, "Failed to get type", http.StatusInternalServerError)
		return
	}

	if t == nil {
		log.Info("Type not found")
		writeJson(w, struct{}{})
	} else {
		log.Info("Type retrieved successfully")
		writeJson(w, t)
	}
}

func HandleGetTypeByName(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleGetTypeByName")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		log.Warn("Missing type name in query")
		http.Error(w, "Missing type name", http.StatusBadRequest)
		return
	}

	log = log.WithField("name", name)
	t, err := st.GetTypeByName(name)
	if err != nil {
		log.WithError(err).Error("Failed to get type by name from DB")
		http.Error(w, "Failed to get type", http.StatusInternalServerError)
		return
	}

	if t == nil {
		log.Info("Type not found by name")
		writeJson(w, struct{}{})
	} else {
		log.Info("Type retrieved successfully by name")
		writeJson(w, t)
	}
}

func HandleCreateType(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleCreateType")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var animalType storage.Types
	if err := jsonDeserialize(r.Body, &animalType); err != nil {
		log.WithError(err).Error("Failed to deserialize animal type body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := st.CreateType(animalType)
	if err != nil {
		log.WithError(err).WithField("animalType", animalType).Error("Failed to create type in DB")
		http.Error(w, "Failed to create type", http.StatusInternalServerError)
		return
	}

	animalType.Id = id
	log.WithFields(logrus.Fields{
		"id":   id,
		"name": animalType.Type,
	}).Info("Animal type created successfully")

	writeJson(w, animalType)
}
