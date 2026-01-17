package store

import (
	"context"
	"database/sql"
)

type ReceiptItem struct {
	ID int64 `json:"id"`

	ReceiptID int64 `json:"receipt_id"`
	ItemID    int64 `json:"item_id"`
	ItemName  string `json:"item_name"`

	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
	Qty           *float64 `json:"qty"`
	Rate          *string  `json:"rate"`

	Total float64 `json:"total"`
	Status int    `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ReceiptItemStore struct {
	db *sql.DB
}

func NewReceiptItemStore(db *sql.DB) *ReceiptItemStore {
	return &ReceiptItemStore{db: db}
}

func (s *ReceiptItemStore) GetByReceiptID(
	ctx context.Context,
	receiptID int64,
) ([]ReceiptItem, error) {

	query := `
		SELECT
			id, receipt_id, item_id, item_name,
			previous_value, current_value, qty, rate,
			total, status, created_at, updated_at
		FROM receipt_items
		WHERE receipt_id = ? AND status = '1'
		ORDER BY id ASC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ReceiptItem

	for rows.Next() {
		var i ReceiptItem
		if err := rows.Scan(
			&i.ID,
			&i.ReceiptID,
			&i.ItemID,
			&i.ItemName,
			&i.PreviousValue,
			&i.CurrentValue,
			&i.Qty,
			&i.Rate,
			&i.Total,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

func (s *ReceiptItemStore) Create(
	ctx context.Context,
	tx *sql.Tx,
	i *ReceiptItem,
) (*int64, error) {

	query := `
		INSERT INTO receipt_items
		(receipt_id, item_id, item_name,
		 previous_value, current_value, qty, rate,
		 total, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, '1')
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		i.ReceiptID,
		i.ItemID,
		i.ItemName,
		i.PreviousValue,
		i.CurrentValue,
		i.Qty,
		i.Rate,
		i.Total,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	i.ID = id
	return &id, nil
}

func (s *ReceiptItemStore) Update(
	ctx context.Context,
	tx *sql.Tx,
	i *ReceiptItem,
) (*int64, error) {

	query := `
		UPDATE receipt_items
		SET item_id = ?, item_name = ?,
		    previous_value = ?, current_value = ?, qty = ?, rate = ?,
		    total = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		i.ItemID,
		i.ItemName,
		i.PreviousValue,
		i.CurrentValue,
		i.Qty,
		i.Rate,
		i.Total,
		i.ID,
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

	return &i.ID, nil
}

func (s *ReceiptItemStore) DeleteByReceiptID(
	ctx context.Context,
	tx *sql.Tx,
	receiptID int64,
) error {

	query := `
		DELETE FROM receipt_items
		WHERE receipt_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, receiptID)
	return err
}

func (s *ReceiptItemStore) Delete(
	ctx context.Context,
	id int64,
) error {

	query := `DELETE FROM receipt_items WHERE id = ?`

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
