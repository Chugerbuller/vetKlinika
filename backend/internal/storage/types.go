package storage

import "database/sql"

func (s *Storage) GetTypes() (*[]Types, error) {
	var types []Types

	query := `SELECT * FROM types;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Types
		err = rows.Scan(&t.Id, &t.Type)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return &types, nil
}
func (s *Storage) GetTypeById(id int) (*Types, error) {
	var t Types

	query := `SELECT * FROM types WHERE id = ?;`
	err := s.db.QueryRow(query, id).Scan(&t.Id, &t.Type)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
func (s *Storage) GetTypeByName(name string) (*Types, error) {
	var t Types

	query := `SELECT * FROM types WHERE type = ?;`
	err := s.db.QueryRow(query, name).Scan(&t.Id, &t.Type)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
func (s *Storage) CreateType(animalType Types) (int, error) {
	query := `INSERT INTO types (type) VALUES (:type);`
	result, err := s.db.Exec(query, sql.Named("type", animalType.Type))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
