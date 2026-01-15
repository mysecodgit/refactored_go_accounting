package store

import (
	"context"
	"database/sql"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID                int64     `json:"id"`
	Type              string    `json:"type"` // invoice | payment | check | deposit | bill
	TransactionDate   string `json:"transaction_date"`
	TransactionNumber string    `json:"transaction_number"`
	Memo              string    `json:"memo"`
	Status            string    `json:"status"` // 0 | 1
	BuildingID        int64     `json:"building_id"`
	UserID            int64     `json:"user_id"`
	UnitID            *int64 `json:"unit_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TransactionStore struct {
	db *sql.DB
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{db: db}
}

// GetAll returns all transactions for a building
func (s *TransactionStore) GetAll(ctx context.Context, buildingID int64) ([]Transaction, error) {
	query := `
		SELECT id, type, transaction_date, transaction_number, memo, status,
		       building_id, user_id, unit_id, created_at, updated_at
		FROM transactions
		WHERE building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID,
			&t.Type,
			&t.TransactionDate,
			&t.TransactionNumber,
			&t.Memo,
			&t.Status,
			&t.BuildingID,
			&t.UserID,
			&t.UnitID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetByID returns a single transaction by ID
func (s *TransactionStore) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	query := `
		SELECT id, type, transaction_date, transaction_number, memo, status,
		       building_id, user_id, unit_id, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var t Transaction
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Type,
		&t.TransactionDate,
		&t.TransactionNumber,
		&t.Memo,
		&t.Status,
		&t.BuildingID,
		&t.UserID,
		&t.UnitID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

// Create inserts a new transaction
func (s *TransactionStore) Create(ctx context.Context, tx *sql.Tx, t *Transaction) (*int64, error) {
	query := `
		INSERT INTO transactions
		(type, transaction_date, transaction_number, memo, status,
		 building_id, user_id, unit_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		t.Type,
		t.TransactionDate,
		t.TransactionNumber,
		t.Memo,
		"1",
		t.BuildingID,
		t.UserID,
		t.UnitID,
	)
	if err != nil {
		return nil, err
	}	

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return &id, nil
}

// Update modifies an existing transaction
func (s *TransactionStore) Update(ctx context.Context, tx *sql.Tx, t *Transaction) (*int64, error) {
	query := `
		UPDATE transactions
		SET type = ?, transaction_date = ?, transaction_number = ?, memo = ?, status = ?,
		    building_id = ?, user_id = ?, unit_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		t.Type,
		t.TransactionDate,
		t.TransactionNumber,
		t.Memo,
		t.Status,
		t.BuildingID,
		t.UserID,
		t.UnitID,
		t.ID,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, ErrNotFound
	}

	return &t.ID, nil
}

// Delete removes a transaction by ID
func (s *TransactionStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
