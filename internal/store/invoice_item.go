package store

import (
	"context"
	"database/sql"
)

type InvoiceItem struct {
	ID int64 `json:"id"`

	InvoiceID int64  `json:"invoice_id"`
	ItemID    int    `json:"item_id"`
	ItemName  string `json:"item_name"`

	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
	Qty           float64  `json:"qty"`
	Rate          float64  `json:"rate"`

	Total  float64 `json:"total"`
	Status *int    `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InvoiceItemStore struct {
	db *sql.DB
}

func NewInvoiceItemStore(db *sql.DB) *InvoiceItemStore {
	return &InvoiceItemStore{db: db}
}

func (s *InvoiceItemStore) GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]InvoiceItem, error) {
	query := `
		SELECT id, invoice_id, item_id, item_name,
		       previous_value, current_value, qty, rate,
		       total, status, created_at, updated_at
		FROM invoice_items
		WHERE invoice_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var i InvoiceItem
		if err := rows.Scan(
			&i.ID,
			&i.InvoiceID,
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

func (s *InvoiceItemStore) GetByID(ctx context.Context, id int64) (*InvoiceItem, error) {
	query := `
		SELECT id, invoice_id, item_id, item_name,
		       previous_value, current_value, qty, rate,
		       total, status, created_at, updated_at
		FROM invoice_items
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var i InvoiceItem
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID,
		&i.InvoiceID,
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &i, nil
}

func (s *InvoiceItemStore) Create(ctx context.Context, tx *sql.Tx, i *InvoiceItem) error {
	query := `
		INSERT INTO invoice_items
		(invoice_id, item_id, item_name,
		 previous_value, current_value, qty, rate,
		 total, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		i.InvoiceID,
		i.ItemID,
		i.ItemName,
		i.PreviousValue,
		i.CurrentValue,
		i.Qty,
		i.Rate,
		i.Total,
		"1",
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	i.ID = id
	return nil
}

func (s *InvoiceItemStore) Update(ctx context.Context, i *InvoiceItem) error {
	query := `
		UPDATE invoice_items
		SET invoice_id = ?, item_id = ?, item_name = ?,
		    previous_value = ?, current_value = ?, qty = ?, rate = ?,
		    total = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		i.InvoiceID,
		i.ItemID,
		i.ItemName,
		i.PreviousValue,
		i.CurrentValue,
		i.Qty,
		i.Rate,
		i.Total,
		i.Status,
		i.ID,
	)
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

func (s *InvoiceItemStore) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE invoice_items SET status = '0' WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, id)
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

func (s *InvoiceItemStore) DeleteByInvoiceID(ctx context.Context, tx *sql.Tx, invoiceID int64) error {
	query := `UPDATE invoice_items SET status = '0' WHERE invoice_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, invoiceID)
	return err
}
