package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"vetklinika/backend/internal/service"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

func HandlePets(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r)
	log.Info("Incoming request to HandlePets")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("ownerId") {
			handleGetPetsByOwnerId(w, r)
		} else if r.URL.Query().Has("id") {
			handleGetPetById(w, r)
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
		log.WithField("method", r.Method).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetPetById(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetPetById")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		log.Warn("Missing pet ID in query")
		http.Error(w, "Missing pet ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Warn("Invalid pet ID format")
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	res, err := st.GetPetById(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Error retrieving pet from DB")
		http.Error(w, "Error retrieving pet", http.StatusInternalServerError)
		return
	}

	log.WithField("id", id).Info("Pet retrieved successfully")
	writeJson(w, res)
}

func handleDeletePet(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleDeletePet")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	if len(strId) == 0 {
		log.Warn("Missing pet ID in query")
		http.Error(w, "Missing pet ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Error("Failed to parse pet ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = st.DeletePet(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Error deleting pet from DB")
		http.Error(w, "Error deleting pet", http.StatusInternalServerError)
		return
	}

	log.WithField("id", id).Info("Pet deleted successfully")
	w.WriteHeader(http.StatusOK)
	writeJson(w, id)
}

func handleGetPetsByOwnerId(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetPetsByOwnerId")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("ownerId")
	if len(strId) == 0 {
		log.Warn("Missing owner ID in query")
		http.Error(w, "Missing owner ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).WithField("strId", strId).Error("Failed to parse owner ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := st.GetPetsByOwner(id)
	if err != nil {
		log.WithError(err).WithField("ownerId", id).Error("Error retrieving pets for owner from DB")
		http.Error(w, "Error retrieving pets", http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{"ownerId": id, "count": len(*res)}).Info("Pets by owner retrieved successfully")
	writeJson(w, res)
}

func handleGetPets(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetPets")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	res, err := st.GetPets()
	if err != nil {
		log.WithError(err).Error("Error retrieving all pets from DB")
		http.Error(w, "Error retrieving pets", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(*res)).Info("All pets retrieved successfully")
	writeJson(w, res)
}

func handleCreatePet(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleCreatePet")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	ownerIdstr := r.URL.Query().Get("ownerId")
	if len(ownerIdstr) == 0 {
		log.Warn("Missing owner ID in query")
		http.Error(w, "Missing owner ID", http.StatusBadRequest)
		return
	}

	ownerId, err := strconv.Atoi(ownerIdstr)
	if err != nil {
		log.WithError(err).WithField("ownerIdstr", ownerIdstr).Warn("Invalid owner ID format")
		http.Error(w, "Invalid owner ID", http.StatusBadRequest)
		return
	}

	type PetRes struct {
		Id           int     `json:"id"`
		Name         string  `json:"name"`
		Weight       float64 `json:"weight"`
		Age          int     `json:"age"`
		AgeWriteDate string  `json:"age_write_date"`
		Type         int     `json:"type"`
		Allergies    string  `json:"allergies"`
	}

	var pet PetRes
	err = jsonDeserialize(r.Body, &pet)
	if err != nil {
		log.WithError(err).Error("Error deserializing pet data")
		http.Error(w, "Error deserializing pet data", http.StatusBadRequest)
		return
	}

	pet.AgeWriteDate, err = service.ConvertDate(pet.AgeWriteDate)
	if err != nil {
		log.WithError(err).WithField("date", pet.AgeWriteDate).Error("Error converting date")
		http.Error(w, "Error creating pet", http.StatusInternalServerError)
		return
	}

	petForSave := storage.Pet{
		Name:         pet.Name,
		Weight:       pet.Weight,
		Age:          pet.Age,
		AgeWriteDate: pet.AgeWriteDate,
		Type:         strconv.Itoa(pet.Type),
		Allergies:    pet.Allergies,
	}

	id, err := st.AddPet(petForSave, ownerId)
	if err != nil || id == 0 {
		log.WithError(err).WithField("ownerId", ownerId).Error("Error adding pet to DB")
		http.Error(w, "Error adding pet", http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{"new_pet_id": id, "ownerId": ownerId}).Info("Pet created successfully")
	w.WriteHeader(http.StatusCreated)
	writeJson(w, id)
}

func handleUpdatePet(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleUpdatePet")
	st := getStorage(r)
	if st == nil {
		log.Error("Storage not available")
		writeError(w, fmt.Errorf("Storage not available"), http.StatusInternalServerError)
		return
	}

	var pet storage.Pet
	err := jsonDeserialize(r.Body, &pet)
	if err != nil {
		log.WithError(err).Error("Error deserializing pet data for update")
		writeError(w, fmt.Errorf("Error deserializing pet data: %v", err), http.StatusBadRequest)
		return
	}

	log.WithField("pet", pet).Debug("Received pet for update")

	err = st.UpdatePet(pet)
	if err != nil {
		log.WithError(err).WithField("id", pet.Id).Error("Error updating pet in DB")
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.WithField("id", pet.Id).Info("Pet updated successfully")
	writeJson(w, pet.Id)
}
