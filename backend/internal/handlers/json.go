package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	writeJson(w, errorResponse{
		Error: err.Error(),
	})
}
func writeJson(w http.ResponseWriter, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, err := w.Write([]byte(response)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func jsonDeserialize(stream io.ReadCloser, dest any) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return err
	}
	defer stream.Close()
	if err := json.Unmarshal(buf.Bytes(), &dest); err != nil {
		return err
	}
	return nil
}
