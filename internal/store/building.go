package store

import (
	"context"
	"database/sql"
	"time"
)

type Building struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BuildingStore struct {
	db *sql.DB
}

func (s *BuildingStore) GetAll(ctx context.Context) ([]Building, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM buildings
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []Building
	for rows.Next() {
		var b Building
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}

	return buildings, nil
}

func (s *BuildingStore) GetByID(ctx context.Context, id int64) (*Building, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM buildings
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var b Building
	err := s.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &b, nil
}

func (s *BuildingStore) Create(ctx context.Context, building *Building) error {
	query := `
		INSERT INTO buildings (name)
		VALUES (?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, building.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	building.ID = id
	return nil
}

func (s *BuildingStore) Update(ctx context.Context, building *Building) error {
	query := `
		UPDATE buildings
		SET name = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, building.Name, building.ID)
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

func (s *BuildingStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM buildings
		WHERE id = ?
	`

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
