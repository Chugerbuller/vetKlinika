package storage

import (
	"database/sql"
)

func (s *Storage) GetAppointmentsByPetId(petId int) (*[]Appointment, error) {
	var appointments []Appointment
	query := `SELECT *
				FROM appointments
				WHERE pet_id = :petId;`
	rows, err := s.db.Query(query, sql.Named("petId", petId))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var appointment Appointment
		err = rows.Scan(&appointment.Id, &appointment.PetId, &appointment.Date, &appointment.Anamnesis, &appointment.Disease, &appointment.Treatment, &appointment.Recommendations)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return &appointments, nil
}
func (s *Storage) GetAppointmentById(id int) (*Appointment, error) {
	var appointment Appointment
	query := `SELECT *
				FROM appointments
				WHERE id = :id;`
	err := s.db.QueryRow(query, sql.Named("id", id)).Scan(&appointment.Id, &appointment.PetId, &appointment.Date, &appointment.Anamnesis, &appointment.Disease, &appointment.Treatment, &appointment.Recommendations)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}
func (s *Storage) CreateAppointment(appointment *Appointment) (int, error) {
	query := `INSERT INTO appointments (pet_id, date, anamnesis, disease, treatment, recommendations)
				VALUES (:petId, :date, :anamnesis, :disease, :treatment, :recommendations);`
	result, err := s.db.Exec(query, sql.Named("petId", appointment.PetId), sql.Named("date", appointment.Date), sql.Named("anamnesis", appointment.Anamnesis), sql.Named("disease", appointment.Disease), sql.Named("treatment", appointment.Treatment), sql.Named("recommendations", appointment.Recommendations))
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}
