package store

import (
	"context"
	"database/sql"
	"time"
)

type Role struct {
	ID          int64     `json:"id"`
	OwnerUserID int64     `json:"owner_user_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetAllByOwnerID(ctx context.Context, ownerUserID int64) ([]Role, error) {
	query := `
		SELECT id, owner_user_id, name, created_at
		FROM roles
		WHERE owner_user_id = ?
		ORDER BY name
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.OwnerUserID, &r.Name, &r.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

func (s *RoleStore) GetByID(ctx context.Context, id int64) (*Role, error) {
	query := `
		SELECT id, owner_user_id, name, created_at
		FROM roles
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var r Role
	err := s.db.QueryRowContext(ctx, query, id).Scan(&r.ID, &r.OwnerUserID, &r.Name, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &r, nil
}

func (s *RoleStore) Create(ctx context.Context, role *Role) error {
	query := `
		INSERT INTO roles (owner_user_id, name)
		VALUES (?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, role.OwnerUserID, role.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	role.ID = id
	return nil
}

func (s *RoleStore) Update(ctx context.Context, role *Role) error {
	query := `
		UPDATE roles
		SET name = ?
		WHERE id = ? AND owner_user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, role.Name, role.ID, role.OwnerUserID)
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

func (s *RoleStore) Delete(ctx context.Context, id int64, ownerUserID int64) error {
	query := `
		DELETE FROM roles
		WHERE id = ? AND owner_user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id, ownerUserID)
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
