package store

import (
	"context"
	"database/sql"
)

type InvoiceAppliedDiscount struct {
	ID int64 `json:"id"`

	Reference     string `json:"reference"`
	InvoiceID     int64  `json:"invoice_id"`
	TransactionID int64  `json:"transaction_id"`

	ARAccountID     int64 `json:"ar_account"`
	IncomeAccountID int64 `json:"income_account"`

	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date         string `json:"date"`

	Status string `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InvoiceAppliedDiscountStore struct {
	db *sql.DB
}

func NewInvoiceAppliedDiscountStore(db *sql.DB) *InvoiceAppliedDiscountStore {
	return &InvoiceAppliedDiscountStore{db: db}
}

func (s *InvoiceAppliedDiscountStore) GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]InvoiceAppliedDiscount, error) {
	query := `
		SELECT id, reference, invoice_id, transaction_id,
		       ar_account, income_account,
		       amount, description, date, status,
		       created_at, updated_at
		FROM invoice_applied_discounts
		WHERE invoice_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []InvoiceAppliedDiscount
	for rows.Next() {
		var d InvoiceAppliedDiscount
		if err := rows.Scan(
			&d.ID,
			&d.Reference,
			&d.InvoiceID,
			&d.TransactionID,
			&d.ARAccountID,
			&d.IncomeAccountID,
			&d.Amount,
			&d.Description,
			&d.Date,
			&d.Status,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		discounts = append(discounts, d)
	}

	return discounts, nil
}

func (s *InvoiceAppliedDiscountStore) GetByID(ctx context.Context, id int64) (*InvoiceAppliedDiscount, error) {
	query := `
		SELECT id, reference, invoice_id, transaction_id,
		       ar_account, income_account,
		       amount, description, date, status,
		       created_at, updated_at
		FROM invoice_applied_discounts
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var d InvoiceAppliedDiscount
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Reference,
		&d.InvoiceID,
		&d.TransactionID,
		&d.ARAccountID,
		&d.IncomeAccountID,
		&d.Amount,
		&d.Description,
		&d.Date,
		&d.Status,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &d, nil
}

func (s *InvoiceAppliedDiscountStore) Create(ctx context.Context, tx *sql.Tx, d *InvoiceAppliedDiscount) error {
	query := `
		INSERT INTO invoice_applied_discounts
		(reference, invoice_id, transaction_id,
		 ar_account, income_account,
		 amount, description, date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		d.Reference,
		d.InvoiceID,
		d.TransactionID,
		d.ARAccountID,
		d.IncomeAccountID,
		d.Amount,
		d.Description,
		d.Date,
		d.Status,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	d.ID = id
	return nil
}

func (s *InvoiceAppliedDiscountStore) Update(ctx context.Context, d *InvoiceAppliedDiscount) error {
	query := `
		UPDATE invoice_applied_discounts
		SET reference = ?, invoice_id = ?, transaction_id = ?,
		    ar_account = ?, income_account = ?,
		    amount = ?, description = ?, date = ?, status = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		d.Reference,
		d.InvoiceID,
		d.TransactionID,
		d.ARAccountID,
		d.IncomeAccountID,
		d.Amount,
		d.Description,
		d.Date,
		d.Status,
		d.ID,
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

func (s *InvoiceAppliedDiscountStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM invoice_applied_discounts WHERE id = ?`

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

func (s *InvoiceAppliedDiscountStore) DeleteByInvoiceID(ctx context.Context, invoiceID int64) error {
	query := `DELETE FROM invoice_applied_discounts WHERE invoice_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, invoiceID)
	return err
}
