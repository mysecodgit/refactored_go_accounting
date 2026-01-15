package store

import (
	"context"
	"database/sql"
	"time"
)

type People struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	TypeID     int64      `json:"type_id"`
	BuildingID int64      `json:"building_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Type       PeopleType `json:"type"`
}

type PeopleStore struct {
	db *sql.DB
}

func (s *PeopleStore) GetAll(ctx context.Context, buildingID int64) ([]People, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			p.phone,
			pt.id,
			pt.title,
			p.created_at,
			p.updated_at
		FROM people p
		JOIN people_types pt ON pt.id = p.type_id
		WHERE p.building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []People

	for rows.Next() {
		var p People
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Phone,
			&p.Type.ID,
			&p.Type.Title,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}

func (s *PeopleStore) GetByID(ctx context.Context, id int64) (*People, error) {
	query := `SELECT id, name, phone, type_id, building_id, created_at, updated_at FROM people WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p People
	err := s.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Phone, &p.TypeID, &p.BuildingID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (s *PeopleStore) Create(ctx context.Context, p *People) error {
	query := `INSERT INTO people (name, phone, type_id, building_id) VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Phone, p.TypeID, p.BuildingID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (s *PeopleStore) Update(ctx context.Context, p *People) error {
	query := `UPDATE people SET name = ?, phone = ?, type_id = ?, building_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Phone, p.TypeID, p.BuildingID, p.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PeopleStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM people WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
