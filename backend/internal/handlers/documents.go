package handlers

import (
	"net/http"
	"strconv"
	"vetklinika/backend/internal/service"
	"vetklinika/backend/internal/storage"
)

func HandleDocuments(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r) // Ваша вспомогательная функция
	log.Info("Incoming request to HandleDocuments")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			HandleGetDocument(w, r)
		} else if r.URL.Query().Has("appointmentId") {
			HandleGetDocumentsByAppointment(w, r)
		} else if r.URL.Query().Has("petId") {
			HandleGetDocumentByPet(w, r)
		} else {
			log.Warn("Missing appointmentId or petId query parameter")
			http.Error(w, "Missing appointmentId or petId query parameter", http.StatusBadRequest)
		}
	case http.MethodPost:
		HandleCreateDocument(w, r)
	case http.MethodDelete:
		HandleDeleteDocument(w, r)
	default:
		log.WithField("method", r.Method).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleDeleteDocument")
	s := getStorage(r) // Ваша вспомогательная функция
	if s == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	docIdStr := r.URL.Query().Get("id")
	if len(docIdStr) == 0 {
		log.Warn("Missing id query parameter")
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}
	docId, err := strconv.Atoi(docIdStr)
	if err != nil {
		log.WithError(err).WithField("docIdStr", docIdStr).Warn("Invalid id query parameter")
		http.Error(w, "Invalid id query parameter", http.StatusBadRequest)
		return
	}
	err = s.DeleteDocument(docId)
	if err != nil {
		log.WithError(err).WithField("docId", docId).Error("Failed to delete document")
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}
	log.WithField("docId", docId).Info("Document deleted successfully")
	writeJson(w, struct{}{})
}

func HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleGetDocument")
	s := getStorage(r)
	if s == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	docIdStr := r.URL.Query().Get("id")
	if len(docIdStr) == 0 {
		log.Warn("Missing id query parameter")
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}
	docId, err := strconv.Atoi(docIdStr)
	if err != nil {
		log.WithError(err).WithField("docIdStr", docIdStr).Warn("Invalid id query parameter")
		http.Error(w, "Invalid id query parameter", http.StatusBadRequest)
		return
	}
	doc, err := s.GetDocumentById(docId)
	if err != nil {
		log.WithError(err).WithField("docId", docId).Error("Failed to get document")
		http.Error(w, "Failed to get document", http.StatusInternalServerError)
		return
	}
	if doc == nil {
		log.WithField("docId", docId).Info("Document not found, returning empty struct")
		writeJson(w, struct{}{})
	} else {
		log.WithField("docId", docId).Info("Document retrieved successfully")
		writeJson(w, doc)
	}
}

func HandleCreateDocument(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleCreateDocument")
	var err error
	s := getStorage(r)
	if s == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var doc storage.Document
	if err := jsonDeserialize(r.Body, &doc); err != nil {
		log.WithError(err).Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	doc.Date, err = service.ConvertDate(doc.Date)
	if err != nil {
		log.WithError(err).WithField("date", doc.Date).Error("Error converting date")
		http.Error(w, "Error creating appointment", http.StatusInternalServerError)
		return
	}
	docId, err := s.CreateDocument(doc)
	if err != nil {
		log.WithError(err).Error("Failed to create document")
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}
	log.WithField("docId", docId).Info("Document created successfully")
	writeJson(w, docId)
}

func HandleGetDocumentByPet(w http.ResponseWriter, r *http.Request) {
	// Оставляю пустым, как в оригинале, только лог вызова
	log := getLogger(r).WithField("handler", "HandleGetDocumentByPet")
	log.Info("HandleGetDocumentByPet called")
}

func HandleGetDocumentsByAppointment(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).WithField("handler", "HandleGetDocumentsByAppointment")
	s := getStorage(r)
	if s == nil {
		log.Error("Storage not available in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	appointmentIdStr := r.URL.Query().Get("appointmentId")
	if len(appointmentIdStr) == 0 {
		log.Warn("Missing appointmentId query parameter")
		http.Error(w, "Missing appointmentId query parameter", http.StatusBadRequest)
		return
	}
	appointmentId, err := strconv.Atoi(appointmentIdStr)
	if err != nil {
		log.WithError(err).WithField("appointmentIdStr", appointmentIdStr).Warn("Invalid appointmentId query parameter")
		http.Error(w, "Invalid appointmentId query parameter", http.StatusBadRequest)
		return
	}
	documents, err := s.GetDocumentsByAppointment(appointmentId)
	if err != nil {
		log.WithError(err).WithField("appointmentId", appointmentId).Error("Failed to get documents")
		http.Error(w, "Failed to get documents"+err.Error(), http.StatusInternalServerError)
		return
	}
	if documents == nil {
		log.WithField("appointmentId", appointmentId).Info("No documents found for appointment")
		writeJson(w, struct{}{})
	} else {
		log.WithField("count", len(documents)).Info("Documents retrieved successfully")
		writeJson(w, documents)
	}
}
