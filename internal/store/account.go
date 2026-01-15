package store

import (
	"context"
	"database/sql"
	"time"
)

type Account struct {
	ID            int64     `json:"id"`
	AccountNumber int       `json:"account_number"`
	AccountName   string    `json:"account_name"`
	AccountType   int64     `json:"account_type"`  // FK → account_types.id
	BuildingID    int64     `json:"building_id"`   // FK → buildings.id
	Type      AccountType  `json:"type"`
	IsDefault     int      `json:"isDefault"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AccountStore struct {
	db *sql.DB
}

func (s *AccountStore) GetAll(ctx context.Context, buildingID int64) ([]Account, error) {
	query := `
		SELECT acc.id, acc.account_number, acc.account_name,at.id,at.typeName, acc.isDefault, acc.created_at, acc.updated_at
		FROM accounts acc
		JOIN account_types at ON acc.account_type = at.id
		WHERE building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(
			&a.ID, 
			&a.AccountNumber, 
			&a.AccountName, 
			&a.Type.ID,
			&a.Type.TypeName,
			&a.IsDefault, 
			&a.CreatedAt, 
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *AccountStore) GetByID(ctx context.Context, id int64) (*Account, error) {
	query := `
		SELECT id, account_number, account_name, account_type, building_id, isDefault, created_at, updated_at
		FROM accounts
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var a Account
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.AccountNumber, &a.AccountName, &a.AccountType, &a.BuildingID, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &a, nil
}

func (s *AccountStore) Create(ctx context.Context, a *Account) error {
	query := `
		INSERT INTO accounts (account_number, account_name, account_type, building_id, isDefault)
		VALUES (?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, a.AccountNumber, a.AccountName, a.AccountType, a.BuildingID, a.IsDefault)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = id
	return nil
}

func (s *AccountStore) Update(ctx context.Context, a *Account) error {
	query := `
		UPDATE accounts
		SET account_number = ?, account_name = ?, account_type = ?, building_id = ?, isDefault = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, a.AccountNumber, a.AccountName, a.AccountType, a.BuildingID, a.IsDefault, a.ID)
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

func (s *AccountStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM accounts WHERE id = ?`

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
