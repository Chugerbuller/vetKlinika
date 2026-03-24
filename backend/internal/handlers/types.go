package handlers

import (
	"log"
	"net/http"
	"strconv"
	"vetklinika/backend/internal/storage"
)

func HandleTypes(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HandleGetTypes(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	types, err := st.GetTypes()
	if err != nil {
		http.Error(w, "Failed to get types", http.StatusInternalServerError)
		return
	}
	if types == nil {
		writeJson(w, struct{}{})
	} else {
		writeJson(w, types)
	}
}
func HandleGetTypeById(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		http.Error(w, "Missing type ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}
	t, err := st.GetTypeById(id)
	if err != nil {
		http.Error(w, "Failed to get type", http.StatusInternalServerError)
		return
	}
	if t == nil {
		writeJson(w, struct{}{})
	} else {
		writeJson(w, t)
	}
}
func HandleGetTypeByName(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		http.Error(w, "Missing type name", http.StatusBadRequest)
		return
	}
	t, err := st.GetTypeByName(name)
	if err != nil {
		http.Error(w, "Failed to get type", http.StatusInternalServerError)
		return
	}
	if t == nil {
		writeJson(w, struct{}{})
	} else {
		writeJson(w, t)
	}
}
func HandleCreateType(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var animalType storage.Types
	if err := jsonDeserialize(r.Body, &animalType); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := st.CreateType(animalType)
	if err != nil {
		http.Error(w, "Failed to create type", http.StatusInternalServerError)
		return
	}
	animalType.Id = id
	log.Printf("Created type %v:", animalType)
	writeJson(w, animalType)

}
