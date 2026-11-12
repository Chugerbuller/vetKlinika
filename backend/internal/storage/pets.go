package storage

import (
	"database/sql"
	"log"
	"vetklinika/backend/internal/service"
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
	defer rows.Close()
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
		if err != nil {
			return nil, err
		}
		pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}
func (s *Storage) GetPetById(id int) (*Pet, error) {
	var pet Pet
	query := `SELECT *
			FROM pets_with_types
			WHERE pet_id = :id;`
	row := s.db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
	if err != nil {
		log.Printf("Error scanning pet with ID %d: %v", id, err)
		return nil, err
	}
	pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
	if err != nil {
		return nil, err
	}
	return &pet, nil
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
	defer rows.Close()
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
		if err != nil {
			return nil, err
		}
		pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
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
	defer rows.Close()
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
		if err != nil {
			return nil, err
		}
		pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
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
			p.allergies

		FROM persons_to_pets ptp
		JOIN pets p ON ptp.pet_id = p.id
		LEFT JOIN types t ON p.type = t.id
		WHERE ptp.person_id = :ownerId;`

	rows, err := s.db.Query(query, sql.Named("ownerId", ownerId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
		if err != nil {
			return nil, err
		}
		pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
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
			p.weight,
			p.age,
			p.age_write_date,
			t.type,
			p.allergies
		FROM persons person
		JOIN persons_to_pets ptp ON person.id = ptp.person_id
		JOIN pets p ON ptp.pet_id = p.id
		LEFT JOIN types t ON p.type = t.id
		WHERE person.name = :ownerName;`

	rows, err := s.db.Query(query, sql.Named("ownerName", ownerName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pet Pet
		err = rows.Scan(&pet.Id, &pet.Name, &pet.Weight, &pet.Age, &pet.AgeWriteDate, &pet.Type, &pet.Allergies)
		if err != nil {
			return nil, err
		}
		pet.Age, err = service.CalcAge(pet.AgeWriteDate, pet.Age)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return &pets, nil
}

// AddPet создаёт питомца и привязывает его к владельцу в одной транзакции.
// Возвращает id созданного питомца.
func (s *Storage) AddPet(pet Pet, personId int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(`INSERT INTO pets(name, type, weight, age, age_write_date, allergies)
			VALUES (:name, :type, :weight, :age, :ageWriteDate, :allergies);`,
		sql.Named("name", pet.Name),
		sql.Named("type", pet.Type),
		sql.Named("weight", pet.Weight),
		sql.Named("age", pet.Age),
		sql.Named("ageWriteDate", pet.AgeWriteDate),
		sql.Named("allergies", pet.Allergies))
	if err != nil {
		return 0, err
	}
	petId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`INSERT INTO persons_to_pets(person_id, pet_id) VALUES (:ownerId, :petId);`,
		sql.Named("ownerId", personId),
		sql.Named("petId", petId))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return int(petId), nil
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
	weight = :weight,
	age = :age,
	age_write_date = :ageWriteDate,
	allergies = :allergies
WHERE id = :id;
`
	log.Printf("PETS: %+v", pet)
	_, err := s.db.Exec(query,
		sql.Named("id", pet.Id),
		sql.Named("name", pet.Name),
		sql.Named("type", pet.Type),
		sql.Named("weight", pet.Weight),
		sql.Named("age", pet.Age),
		sql.Named("ageWriteDate", pet.AgeWriteDate),
		sql.Named("allergies", pet.Allergies))
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
