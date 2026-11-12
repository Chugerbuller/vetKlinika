package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"vetklinika/backend/internal/service"
	"vetklinika/backend/internal/storage"
)

func HandleAppointments(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r)
	log.Info("Incoming request to HandleAppointments")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("petId") {
			handleGetAppointmentsByPetId(w, r)
		} else if r.URL.Query().Has("id") {
			handleGetAppointmentById(w, r)
		} else {
			log.Warn("Missing query parameters in GET request")
			http.Error(w, "Missing petId or id query parameter", http.StatusBadRequest)
		}
	case http.MethodPost:
		handleCreateAppointment(w, r)
	case http.MethodPut:
		handleUpdateAppointment(w, r)
	default:
		log.WithField("method", r.Method).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleCreateAppointment")
	st := getStorage(r)

	if st == nil {
		log.Error("Storage dependency missing in context")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	var appointment storage.Appointment
	err := jsonDeserialize(r.Body, &appointment)
	if err != nil {
		log.WithError(err).Error("Failed to deserialize appointment body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.WithField("appointment_data", appointment).Debug("Attempting to create appointment")

	appointment.Date, err = service.ConvertDate(appointment.Date)
	if err != nil {
		log.WithError(err).WithField("date", appointment.Date).Error("Date conversion failed")
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	id, err := st.CreateAppointment(&appointment)
	if err != nil {
		log.WithError(err).Error("Database error during CreateAppointment")
		http.Error(w, "Error creating appointment", http.StatusInternalServerError)
		return
	}

	log.WithField("created_id", id).Info("Appointment successfully created")
	w.WriteHeader(http.StatusCreated)
	writeJson(w, id)
}

func handleGetAppointmentById(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetAppointmentById")
	st := getStorage(r)

	if st == nil {
		log.Error("Storage missing")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("id")
	log = log.WithField("query_id", strId)

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).Warn("Invalid ID format in query")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	appointment, err := st.GetAppointmentById(id)
	if err != nil {
		log.WithError(err).Error("Failed to fetch appointment from DB")
		http.Error(w, "Error fetching appointment", http.StatusInternalServerError)
		return
	}

	if appointment == nil {
		log.Info("Appointment not found")
		http.NotFound(w, r)
		return
	}

	log.Info("Successfully fetched appointment")
	writeJson(w, appointment)
}

func handleGetAppointmentsByPetId(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetAppointmentsByPetId")
	st := getStorage(r)

	if st == nil {
		log.Error("Storage missing")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}

	strId := r.URL.Query().Get("petId")
	log = log.WithField("pet_id", strId)

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.WithError(err).Warn("Invalid PetID format in query")
		http.Error(w, "Invalid PetID format", http.StatusBadRequest)
		return
	}

	appointments, err := st.GetAppointmentsByPetId(id)
	if err != nil {
		log.WithError(err).Error("Failed to fetch appointments from DB")
		http.Error(w, "Error fetching appointments", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(*appointments)).Info("Successfully fetched appointments for pet")
	writeJson(w, appointments)
}
func handleUpdateAppointment(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "handleGetAppointmentsByPetId")
	st := getStorage(r)

	if st == nil {
		log.Error("Storage missing")
		http.Error(w, "Storage not available", http.StatusInternalServerError)
		return
	}
	var appointment storage.Appointment
	err := jsonDeserialize(r.Body, &appointment)
	if err != nil {
		log.WithError(err).Error("Failed to deserialize appointment body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = st.UpdateAppointment(&appointment)
	if errors.Is(err, sql.ErrNoRows) {
		log.WithError(err).Error(err.Error())
		http.Error(w, "appointment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.WithError(err).Error(err.Error())
		http.Error(w, "Storage problem", http.StatusInternalServerError)
		return
	}
	writeJson(w, struct{}{})
}
