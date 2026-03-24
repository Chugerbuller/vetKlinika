package storage

import (
	"database/sql"
	"log"
	"strings"
)

func (s *Storage) GetPersons() (*[]Person, error) {
	var persons []Person
	query := `SELECT * FROM persons;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var person Person
		err = rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Phone, &person.Address)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &persons, nil
}
func (s *Storage) GetPersonById(id int) (*Person, error) {
	var person Person
	query := `SELECT * FROM persons WHERE id = :id;`

	err := s.db.QueryRow(query, sql.Named("id", id)).Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Phone, &person.Address)
	if err != nil {
		return nil, err
	}

	return &person, nil
}
func (s *Storage) GetPersonsByName(name string) (*[]Person, error) {
	var persons []Person
	log.Printf("Searching for persons with name containing: %s", name)
	query := `SELECT * FROM persons WHERE concat_ws(' ', name, surname, patronymic) LIKE :name;`
	rows, err := s.db.Query(query, sql.Named("name", "%"+name+"%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var person Person
		err = rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Phone, &person.Address)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &persons, nil
}
func (s *Storage) AddPerson(person Person) (int, error) {
	person.Name = strings.ToLower(person.Name)
	person.Surname = strings.ToLower(person.Surname)
	person.Patronymic = strings.ToLower(person.Patronymic)

	query := `INSERT INTO persons (name, surname, patronymic, phone, address) VALUES (:name, :surname, :patronymic, :phone, :address) RETURNING id;`
	var id int
	err := s.db.QueryRow(query, sql.Named("name", person.Name), sql.Named("surname", person.Surname), sql.Named("patronymic", person.Patronymic), sql.Named("phone", person.Phone), sql.Named("address", person.Address)).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (s *Storage) UpdatePerson(person Person) error {
	query := `UPDATE persons SET name = :name, surname = :surname, patronymic = :patronymic, phone = :phone, address = :address WHERE id = :id;`
	_, err := s.db.Exec(query, sql.Named("name", person.Name), sql.Named("surname", person.Surname), sql.Named("patronymic", person.Patronymic), sql.Named("phone", person.Phone), sql.Named("address", person.Address), sql.Named("id", person.Id))
	return err
}
func (s *Storage) DeletePerson(id int) error {
	query := `DELETE FROM persons WHERE id = :id;`
	_, err := s.db.Exec(query, sql.Named("id", id))
	return err
}
