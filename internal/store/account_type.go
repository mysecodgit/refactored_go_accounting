package store

import (
	"context"
	"database/sql"
	"time"
)

type AccountType struct {
	ID        int64     `json:"id"`
	TypeName  string    `json:"typeName"`
	Type      string    `json:"type"`
	SubType   string    `json:"sub_type"`
	TypeStatus string   `json:"typeStatus"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountTypeStore struct {
	db *sql.DB
}

func (s *AccountTypeStore) GetAll(ctx context.Context) ([]AccountType, error) {
	query := `SELECT id, typeName, type, sub_type, typeStatus, created_at, updated_at FROM account_types`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AccountType
	for rows.Next() {
		var at AccountType
		if err := rows.Scan(&at.ID, &at.TypeName, &at.Type, &at.SubType, &at.TypeStatus, &at.CreatedAt, &at.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, at)
	}

	return list, nil
}

func (s *AccountTypeStore) GetByID(ctx context.Context, id int64) (*AccountType, error) {
	query := `SELECT id, typeName, type, sub_type, typeStatus, created_at, updated_at FROM account_types WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var at AccountType
	err := s.db.QueryRowContext(ctx, query, id).Scan(&at.ID, &at.TypeName, &at.Type, &at.SubType, &at.TypeStatus, &at.CreatedAt, &at.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &at, nil
}

func (s *AccountTypeStore) Create(ctx context.Context, at *AccountType) error {
	query := `INSERT INTO account_types (typeName, type, sub_type, typeStatus) VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, at.TypeName, at.Type, at.SubType, at.TypeStatus)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	at.ID = id
	return nil
}

func (s *AccountTypeStore) Update(ctx context.Context, at *AccountType) error {
	query := `UPDATE account_types SET typeName = ?, type = ?, sub_type = ?, typeStatus = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, at.TypeName, at.Type, at.SubType, at.TypeStatus, at.ID)
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

func (s *AccountTypeStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM account_types WHERE id = ?`

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
