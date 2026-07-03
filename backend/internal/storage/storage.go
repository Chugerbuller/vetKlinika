package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to database at %s", dbPath)
	stmt, err := db.Prepare(`
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS pets(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name CHAR(32) NOT NULL DEFAULT "",
    weight REAL NOT NULL DEFAULT 0,
    age INTEGER NOT NULL DEFAULT 0,
    age_write_date CHAR(8) NOT NULL DEFAULT "",
    type INTEGER,
    allergies TEXT NOT NULL DEFAULT "",
    FOREIGN KEY (type) REFERENCES types(id)
);
CREATE TABLE IF NOT EXISTS appointments(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER,
    date CHAR(8) NOT NULL DEFAULT "",
    anamnesis TEXT NOT NULL DEFAULT "",
    disease CHAR(128) NOT NULL DEFAULT "",
    treatment TEXT NOT NULL DEFAULT "",
    recommendations TEXT NOT NULL DEFAULT "",
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);
CREATE TABLE IF NOT EXISTS persons(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name CHAR(32) NOT NULL DEFAULT "",
    surname CHAR(32) NOT NULL DEFAULT "",
    patronymic CHAR(32) NOT NULL DEFAULT "",
    phone CHAR(10) NOT NULL DEFAULT "",
    address CHAR(128) NOT NULL DEFAULT ""
);
CREATE TABLE IF NOT EXISTS  types(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type CHAR(32) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS  persons_to_pets(
    person_id INTEGER,
    pet_id INTEGER
);
CREATE TABLE IF NOT EXISTS documents(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER,
    appointment_id INTEGER,
    date CHAR(8) NOT NULL DEFAULT "",
    name CHAR(32) NOT NULL DEFAULT "",
    content BLOB NOT NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);
CREATE VIEW IF NOT EXISTS pets_with_types AS
SELECT
    pets.id AS pet_id,
    pets.name AS pet_name,
    pets.weight AS pet_weight,
    pets.age AS pet_age,
    pets.age_write_date AS pet_age_write_date,
    types.type AS type_name,
    pets.allergies
FROM pets
LEFT JOIN types ON pets.type = types.id;
CREATE INDEX IF NOT EXISTS pet_name_idx ON pets(name);
CREATE INDEX IF NOT EXISTS appointment_date_idx ON appointments(date);
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}
func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}
