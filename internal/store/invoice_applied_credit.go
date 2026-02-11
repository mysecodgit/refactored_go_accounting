package store

import (
	"context"
	"database/sql"
)

type InvoiceAppliedCredit struct {
	ID int64 `json:"id"`

	InvoiceID    int64 `json:"invoice_id"`
	CreditMemoID int64 `json:"credit_memo_id"`

	Amount      float64 `json:"amount"`
	AmountCents int64   `json:"amount_cents"`
	Description string  `json:"description"`
	Date        string  `json:"date"`

	Status string `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InvoiceAppliedCreditStore struct {
	db *sql.DB
}

func NewInvoiceAppliedCreditStore(db *sql.DB) *InvoiceAppliedCreditStore {
	return &InvoiceAppliedCreditStore{db: db}
}

func (s *InvoiceAppliedCreditStore) GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]InvoiceAppliedCredit, error) {
	query := `
		SELECT id, invoice_id, credit_memo_id,
		       amount, amount_cents, description, date, status,
		       created_at, updated_at
		FROM invoice_applied_credits
		WHERE invoice_id = ? AND status = '1'
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []InvoiceAppliedCredit
	for rows.Next() {
		var c InvoiceAppliedCredit
		if err := rows.Scan(
			&c.ID,
			&c.InvoiceID,
			&c.CreditMemoID,
			&c.Amount,
			&c.AmountCents,
			&c.Description,
			&c.Date,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}

	return credits, nil
}

func (s *InvoiceAppliedCreditStore) GetAllByCreditMemoID(ctx context.Context, creditMemoID int64) ([]InvoiceAppliedCredit, error) {
	query := `
		SELECT id, invoice_id, credit_memo_id,
		       amount, amount_cents, description, date, status,
		       created_at, updated_at
		FROM invoice_applied_credits
		WHERE credit_memo_id = ? AND status = '1'
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, creditMemoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []InvoiceAppliedCredit
	for rows.Next() {
		var c InvoiceAppliedCredit
		if err := rows.Scan(
			&c.ID,
			&c.InvoiceID,
			&c.CreditMemoID,
			&c.Amount,
			&c.AmountCents,
			&c.Description,
			&c.Date,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}

	return credits, nil
}

func (s *InvoiceAppliedCreditStore) GetByID(ctx context.Context, id int64) (*InvoiceAppliedCredit, error) {
	query := `
		SELECT id, invoice_id, credit_memo_id,
		       amount, amount_cents, description, date, status,
		       created_at, updated_at
		FROM invoice_applied_credits
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var c InvoiceAppliedCredit
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.InvoiceID,
		&c.CreditMemoID,
		&c.Amount,
		&c.AmountCents,
		&c.Description,
		&c.Date,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (s *InvoiceAppliedCreditStore) Create(ctx context.Context, c *InvoiceAppliedCredit) error {
	query := `
		INSERT INTO invoice_applied_credits
		(invoice_id, credit_memo_id, amount, amount_cents, description, date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		c.InvoiceID,
		c.CreditMemoID,
		c.Amount,
		c.AmountCents,
		c.Description,
		c.Date,
		c.Status,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = id
	return nil
}

func (s *InvoiceAppliedCreditStore) Update(ctx context.Context, c *InvoiceAppliedCredit) error {
	query := `
		UPDATE invoice_applied_credits
		SET invoice_id = ?, credit_memo_id = ?, amount = ?, amount_cents = ?,
		    description = ?, date = ?, status = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		c.InvoiceID,
		c.CreditMemoID,
		c.Amount,
		c.AmountCents,
		c.Description,
		c.Date,
		c.Status,
		c.ID,
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

func (s *InvoiceAppliedCreditStore) Delete(ctx context.Context, id int64) error {
	query := `UPDATE invoice_applied_credits SET status = '0' WHERE id = ?`

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

func (s *InvoiceAppliedCreditStore) DeleteByInvoiceID(ctx context.Context, invoiceID int64) error {
	query := `UPDATE invoice_applied_credits SET status = '0' WHERE invoice_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, invoiceID)
	return err
}
