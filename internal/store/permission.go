package store

import (
	"context"
	"database/sql"
)

type Permission struct {
	ID     int64  `json:"id"`
	Module string `json:"module"`
	Action string `json:"action"`
	Key    string `json:"key"`
}

type PermissionStore struct {
	db *sql.DB
}

func (s *PermissionStore) GetAll(ctx context.Context) ([]Permission, error) {
	query := `
		SELECT id, module, action, ` + "`key`" + `
		FROM permissions
		ORDER BY module, action
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Module, &p.Action, &p.Key); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

func (s *PermissionStore) GetByID(ctx context.Context, id int64) (*Permission, error) {
	query := `
		SELECT id, module, action, ` + "`key`" + `
		FROM permissions
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p Permission
	err := s.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Module, &p.Action, &p.Key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (s *PermissionStore) GetByKey(ctx context.Context, key string) (*Permission, error) {
	query := `
		SELECT id, module, action, ` + "`key`" + `
		FROM permissions
		WHERE ` + "`key`" + ` = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p Permission
	err := s.db.QueryRowContext(ctx, query, key).Scan(&p.ID, &p.Module, &p.Action, &p.Key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (s *PermissionStore) Create(ctx context.Context, permission *Permission) error {
	query := `
		INSERT INTO permissions (module, action, ` + "`key`" + `)
		VALUES (?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, permission.Module, permission.Action, permission.Key)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	permission.ID = id
	return nil
}

func (s *PermissionStore) Update(ctx context.Context, permission *Permission) error {
	query := `
		UPDATE permissions
		SET module = ?, action = ?, ` + "`key`" + ` = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, permission.Module, permission.Action, permission.Key, permission.ID)
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

func (s *PermissionStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM permissions
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
