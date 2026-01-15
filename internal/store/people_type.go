package store

import (
	"context"
	"database/sql"
)

type PeopleType struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type PeopleTypeStore struct {
	db *sql.DB
}

func (s *PeopleTypeStore) GetAll(ctx context.Context) ([]PeopleType, error) {
	query := `SELECT id, title FROM people_types`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []PeopleType
	for rows.Next() {
		var t PeopleType
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	return types, nil
}

func (s *PeopleTypeStore) GetByID(ctx context.Context, id int64) (*PeopleType, error) {
	query := `SELECT id, title FROM people_types WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var t PeopleType
	err := s.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (s *PeopleTypeStore) Create(ctx context.Context, pt *PeopleType) error {
	query := `INSERT INTO people_types (title) VALUES (?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, pt.Title)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	pt.ID = id
	return nil
}

func (s *PeopleTypeStore) Update(ctx context.Context, pt *PeopleType) error {
	query := `UPDATE people_types SET title = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, pt.Title, pt.ID)
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

func (s *PeopleTypeStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM people_types WHERE id = ?`

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
