package storage

import (
	"database/sql"
)

func (s *Storage) GetDocumentById(id int) (*Document, error) {
	query := `SELECT id, pet_id, date, appointment_id, name, content FROM documents WHERE id = :id`
	var doc Document
	err := s.db.QueryRow(query, sql.Named("id", id)).Scan(&doc.Id, &doc.PetId, &doc.Date, &doc.AppointmentId, &doc.Name, &doc.Content)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
func (s *Storage) GetDocumentsByAppointment(appointmentID int) ([]Document, error) {
	query := `SELECT id, pet_id, appointment_id, date, name, content FROM documents WHERE appointment_id = :appointment_id`
	rows, err := s.db.Query(query, sql.Named("appointment_id", appointmentID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.Id, &doc.PetId, &doc.AppointmentId, &doc.Date, &doc.Name, &doc.Content); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}
func (s *Storage) GetDocumentByPet(petId int) (Document, error) {
	query := `SELECT * FROM documents WHERE pet_id = :pet_id`
	var doc Document
	err := s.db.QueryRow(query, sql.Named("pet_id", doc.PetId)).Scan(&doc.Id, &doc.PetId, &doc.Date, &doc.AppointmentId, &doc.Name, &doc.Content)
	if err != nil {
		return Document{}, err
	}
	return doc, nil
}
func (s *Storage) CreateDocument(doc Document) (int, error) {
	query := `INSERT INTO documents (pet_id, appointment_id,date, name, content) VALUES (:pet_id, :appointment_id, :date, :name, :content)`
	result, err := s.db.Exec(query, sql.Named("pet_id", doc.PetId), sql.Named("appointment_id", doc.AppointmentId), sql.Named("date", doc.Date), sql.Named("name", doc.Name), sql.Named("content", doc.Content))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (s *Storage) DeleteDocument(id int) error {
	query := `DELETE FROM documents WHERE id = :id`
	_, err := s.db.Exec(query, sql.Named("id", id))
	return err
}
