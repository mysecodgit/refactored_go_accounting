package store

import (
	"context"
	"database/sql"
)

type InvoicePayment struct {
	ID int64 `json:"id"`

	TransactionID int64     `json:"transaction_id"`
	Reference     string    `json:"reference"`
	Date          string `json:"date"`

	InvoiceID int64 `json:"invoice_id"`
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`

	Amount float64 `json:"amount"`
	Status string  `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InvoicePaymentStore struct {
	db *sql.DB
}

func NewInvoicePaymentStore(db *sql.DB) *InvoicePaymentStore {
	return &InvoicePaymentStore{db: db}
}

func (s *InvoicePaymentStore) GetAll(ctx context.Context, buildingID int64, startDate *string, endDate *string, peopleID *int, status *string) ([]InvoicePayment, error) {
	query := `
		SELECT ip.id, ip.transaction_id, ip.reference, ip.date, ip.invoice_id, ip.user_id, ip.account_id, ip.amount, ip.status, ip.createdAt, ip.updatedAt 
		FROM invoice_payments ip
		INNER JOIN invoices i ON ip.invoice_id = i.id
		WHERE i.building_id = ?
	`
	args := []interface{}{buildingID}
	
	// Add filters
	if startDate != nil && *startDate != "" {
		query += " AND DATE(ip.date) >= ?"
		args = append(args, *startDate)
	}
	
	if endDate != nil && *endDate != "" {
		query += " AND DATE(ip.date) <= ?"
		args = append(args, *endDate)
	}
	
	if peopleID != nil && *peopleID > 0 {
		query += " AND i.people_id = ?"
		args = append(args, *peopleID)
	}
	
	if status != nil && *status != "" {
		query += " AND ip.status = ?"
		args = append(args, *status)
	}
	
	query += " ORDER BY ip.createdAt DESC"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []InvoicePayment
	for rows.Next() {
		var p InvoicePayment
		if err := rows.Scan(
			&p.ID,
			&p.TransactionID,
			&p.Reference,
			&p.Date,
			&p.InvoiceID,
			&p.UserID,
			&p.AccountID,
			&p.Amount,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (s *InvoicePaymentStore) GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]InvoicePayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       invoice_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM invoice_payments
		WHERE invoice_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []InvoicePayment
	for rows.Next() {
		var p InvoicePayment
		if err := rows.Scan(
			&p.ID,
			&p.TransactionID,
			&p.Reference,
			&p.Date,
			&p.InvoiceID,
			&p.UserID,
			&p.AccountID,
			&p.Amount,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (s *InvoicePaymentStore) GetByID(ctx context.Context, id int64) (*InvoicePayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       invoice_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM invoice_payments
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p InvoicePayment
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.TransactionID,
		&p.Reference,
		&p.Date,
		&p.InvoiceID,
		&p.UserID,
		&p.AccountID,
		&p.Amount,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (s *InvoicePaymentStore) GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*InvoicePayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       invoice_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM invoice_payments
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p InvoicePayment
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.TransactionID,
		&p.Reference,
		&p.Date,
		&p.InvoiceID,
		&p.UserID,
		&p.AccountID,
		&p.Amount,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (s *InvoicePaymentStore) Create(ctx context.Context, tx *sql.Tx, p *InvoicePayment) (*InvoicePayment, error) {
	query := `
		INSERT INTO invoice_payments
		(transaction_id, reference, date,
		 invoice_id, user_id, account_id,
		 amount, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, "1")
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		p.TransactionID,
		p.Reference,
		p.Date,
		p.InvoiceID,
		p.UserID,
		p.AccountID,
		p.Amount,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetByIDTx(ctx, tx, id)
}

func (s *InvoicePaymentStore) Update(ctx context.Context, tx *sql.Tx, p *InvoicePayment) (*InvoicePayment, error) {
	query := `
		UPDATE invoice_payments
		SET transaction_id = ?, reference = ?, date = ?,
		    invoice_id = ?, user_id = ?, account_id = ?,
		    amount = ?, status = ?, updatedAt = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		p.TransactionID,
		p.Reference,
		p.Date,
		p.InvoiceID,
		p.UserID,
		p.AccountID,
		p.Amount,
		p.Status,
		p.ID,
	)
	if err != nil {
		return nil, err
	}
	
	return s.GetByIDTx(ctx, tx, p.ID)
}

func (s *InvoicePaymentStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM invoice_payments WHERE id = ?`

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
