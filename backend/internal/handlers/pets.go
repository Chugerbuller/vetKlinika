package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vetklinika/backend/internal/storage"
)

func HandlePets(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("ownerid") {
			handleGetPetsByOwnerId(w, r)
		} else {
			handleGetPets(w, r)
		}
	case http.MethodPost:
		handleCreatePet(w, r)
	case http.MethodPut:
		handleUpdatePet(w, r)
	case http.MethodDelete:
		handleDeletePet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleDeletePet(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		http.Error(w, "Missing pet ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = st.DeletePet(id)
	if err != nil {
		http.Error(w, "Error deleting pet", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, id)
}
func handleGetPetsByOwnerId(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	strId := r.URL.Query().Get("ownerid")
	if len(strId) == 0 {
		http.Error(w, "Missing owner ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := st.GetPetsByOwner(id)
	if err != nil {
		http.Error(w, "Error retrieving pets", http.StatusInternalServerError)
		return
	}
	writeJson(w, res)
}
func handleGetPets(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	res, err := st.GetPets()

	if err != nil {
		http.Error(w, "Error retrieving pets", http.StatusInternalServerError)
		return
	}
	writeJson(w, res)
}
func handleCreatePet(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	ownerIdstr := r.URL.Query().Get("ownerId")
	if len(ownerIdstr) == 0 {
		http.Error(w, "Missing owner ID", http.StatusBadRequest)
		return
	}
	ownerId, err := strconv.Atoi(ownerIdstr)
	if err != nil {
		http.Error(w, "Invalid owner ID", http.StatusBadRequest)
		return
	}
	var pet storage.Pet
	err = jsonDeserialize(r.Body, &pet)
	if err != nil {
		http.Error(w, "Error deserializing pet data", http.StatusBadRequest)
		return
	}
	id, err := st.AddPet(pet, ownerId)
	if err != nil || id == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJson(w, id)

}
func handleUpdatePet(w http.ResponseWriter, r *http.Request) {
	st, ok := r.Context().Value("storage").(*storage.Storage)
	if !ok {
		writeError(w, fmt.Errorf("Storage not available"), http.StatusInternalServerError)
		return
	}
	var pet storage.Pet
	err := jsonDeserialize(r.Body, &pet)
	log.Printf("Received pet for update: %+v", pet)
	if err != nil {
		writeError(w, fmt.Errorf("Error deserializing pet data: %v", err), http.StatusBadRequest)
		return
	}
	err = st.UpdatePet(pet)
	if err != nil {
		log.Printf("Error updating pet with ID: %d, error: %v", pet.Id, err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Updated pet with ID: %d", pet.Id)
	writeJson(w, pet.Id)
}
