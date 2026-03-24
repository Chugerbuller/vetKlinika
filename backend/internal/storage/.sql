PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS pets(
    id INTEGER AUTOINCREMENT PRIMARY KEY,
    name CHAR(32) NOT NULL DEFAULT "",
    type INTEGER,
    anamnesis TEXT NOT NULL DEFAULT "",
    FOREIGN KEY (type) REFERENCES types(id)
);
CREATE TABLE IF NOT EXISTS appointments(
    id INTEGER AUTOINCREMENT PRIMARY KEY,
    pet_id INTEGER,
    date CHAR(8) NOT NULL DEFAULT "",
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);
CREATE TABLE IF NOT EXISTS persons(
    id INTEGER AUTOINCREMENT PRIMARY KEY,
    name CHAR(32) NOT NULL DEFAULT "",
    surname CHAR(32) NOT NULL DEFAULT "",
    patronymic CHAR(32) NOT NULL DEFAULT "",
    phone CHAR(10) NOT NULL DEFAULT "",
    address CHAR(128) NOT NULL DEFAULT ""
);
CREATE TABLE IF NOT EXISTS  types(
    id INTEGER AUTOINCREMENT PRIMARY KEY,
    type CHAR(32) NOT NULL DEFAULT ""
);
CREATE TABLE IF NOT EXISTS  persons_to_pets(
    person_id INTEGER,
    pet_id INTEGER
);
CREATE TABLE IF NOT EXISTS diagnosis(
    id AUTOINCREMENT PRIMARY KEY,
    disease CHAR(128) NOT NULL DEFAULT "",
    treatment TEXT NOT NULL DEFAULT "",
    recommendations TEXT NOT NULL DEFAULT ""
);
CREATE TABLE IF NOT EXISTS  diagnosis_to_pets(
    diagnosis_id INTEGER,
    pet_id INTEGER
);
CREATE TABLE IF NOT EXISTS  diagnosis_to_appointments(
    diagnosis_id INTEGER,
    appointment_id INTEGER
);
CREATE TABLE IF NOT EXISTS documents(
    id AUTOINCREMENT PRIMARY KEY,
    pet_id INTEGER,
    name CHAR(32) NOT NULL DEFAULT "",
    content BLOB NOT NULL,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);
CREATE INDEX IF NOT EXISTS pet_name_idx ON pets(name);
CREATE INDEX IF NOT EXISTS appointment_date_idx ON appointments(date);


