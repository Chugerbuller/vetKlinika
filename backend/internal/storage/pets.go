package storage

import (
	"database/sql"
	"log"
)

// TODO не работает получение питомца по id владельца, нужно разобраться
func (s *Storage) GetPets() (*[]Pet, error) {
	var pets []Pet
	query := `SELECT *
FROM pets_with_types;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Anamnesis)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}
func (s *Storage) GetPetByName(name string) (*[]Pet, error) {
	var pets []Pet
	query := `SELECT *
FROM pets_with_types
WHERE pet_name = :name;`

	rows, err := s.db.Query(query, sql.Named("name", name))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Anamnesis)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}
func (s *Storage) GetPetsByType(typeName string) (*[]Pet, error) {
	var pets []Pet
	query := `SELECT *
FROM pets_with_types
WHERE type_name = :type;`

	rows, err := s.db.Query(query, sql.Named("type", typeName))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Anamnesis)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}
func (s *Storage) GetPetsByOwner(ownerId int) (*[]Pet, error) {
	log.Printf("OwnerId = %d", ownerId)
	var pets []Pet
	query := `SELECT
			p.id,
			p.name,
			p.weight,
			p.age,
			p.age_write_date,
			t.type,
			p.anamnesis
		FROM persons_to_pets ptp
		JOIN pets p ON ptp.pet_id = p.id
		LEFT JOIN types t ON p.type = t.id
		WHERE ptp.person_id = :ownerId;`

	rows, err := s.db.Query(query, sql.Named("ownerId", ownerId))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Anamnesis)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	log.Printf("Pets: %+v", pets)
	return &pets, nil
}
func (s *Storage) GetPetsByOwnerName(ownerName string) (*[]Pet, error) {
	var pets []Pet
	query := `SELECT
			p.id,
			p.name,
			t.type,
			p.anamnesis
		FROM persons person
		JOIN persons_to_pets ptp ON person.id = ptp.person_id
		JOIN pets pet ON ptp.pet_id = pet.id
		LEFT JOIN types t ON pet.type = t.id
		WHERE person.name = :ownerName;`

	rows, err := s.db.Query(query, sql.Named("ownerName", ownerName))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Anamnesis)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}
func (s *Storage) AddPet(pet Pet, personId int) (int, error) {
	query := `-- Вставляем тип (если его нет)
			INSERT INTO types(type)
			SELECT :type
			WHERE NOT EXISTS (SELECT 1 FROM types WHERE type = :type);

			-- Вставляем питомца
			INSERT INTO pets(name, type, weight, age, age_write_date, anamnesis)
			VALUES (
				:name,
				(SELECT id FROM types WHERE type = :type),
				:weight,
				:age,
				:ageWriteDate,
				:anamnesis
			);
			INSERT INTO persons_to_pets(person_id, pet_id)
			VALUES (
				:ownerId,
				last_insert_rowid()
			);`
	res, err := s.db.Exec(query, sql.Named("name", pet.Name), sql.Named("type", pet.Type), sql.Named("weight", pet.Weight), sql.Named("age", pet.Age), sql.Named("ageWriteDate", pet.AgeWriteDate), sql.Named("anamnesis", pet.Anamnesis), sql.Named("personId", personId))
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}
func (s *Storage) UpdatePet(pet Pet) error {
	query := `
INSERT INTO types(type)
SELECT :type
WHERE NOT EXISTS (SELECT 1 FROM types WHERE type = :type);

UPDATE pets
SET
    name = :name,
    type = (SELECT id FROM types WHERE type = :type),
    anamnesis = :anamnesis
WHERE id = :id;
`
	log.Printf("PETS: %+v", pet)
	_, err := s.db.Exec(query,
		sql.Named("id", pet.Id),
		sql.Named("name", pet.Name),
		sql.Named("type", pet.Type),
		sql.Named("anamnesis", pet.Anamnesis))
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) DeletePet(petId int) error {
	query := `DELETE FROM pets WHERE id=:id;`

	_, err := s.db.Exec(query, sql.Named("id", petId))
	return err
}
