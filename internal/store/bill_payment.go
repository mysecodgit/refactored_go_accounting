package store

import (
	"context"
	"database/sql"
)

type BillPayment struct {
	ID int64 `json:"id"`

	TransactionID int64  `json:"transaction_id"`
	Reference     string `json:"reference"`
	Date          string `json:"date"`

	BillID   int64 `json:"bill_id"`
	UserID   int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`

	Amount float64 `json:"amount"`
	Status string  `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BillPaymentStore struct {
	db *sql.DB
}

func NewBillPaymentStore(db *sql.DB) *BillPaymentStore {
	return &BillPaymentStore{db: db}
}

func (s *BillPaymentStore) GetAll(ctx context.Context, buildingID int64, startDate *string, endDate *string, peopleID *int, status *string) ([]BillPayment, error) {
	query := `
		SELECT bp.id, bp.transaction_id, bp.reference, bp.date,
		       bp.bill_id, bp.user_id, bp.account_id,
		       bp.amount, bp.status, bp.createdAt, bp.updatedAt
		FROM bill_payments bp
		INNER JOIN bills b ON bp.bill_id = b.id
		WHERE b.building_id = ?
	`
	args := []any{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND DATE(bp.date) >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		query += " AND DATE(bp.date) <= ?"
		args = append(args, *endDate)
	}
	if peopleID != nil && *peopleID > 0 {
		query += " AND b.people_id = ?"
		args = append(args, *peopleID)
	}
	if status != nil && *status != "" {
		query += " AND bp.status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY bp.createdAt DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []BillPayment
	for rows.Next() {
		var p BillPayment
		if err := rows.Scan(
			&p.ID,
			&p.TransactionID,
			&p.Reference,
			&p.Date,
			&p.BillID,
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

func (s *BillPaymentStore) GetAllByBillID(ctx context.Context, billID int64) ([]BillPayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       bill_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM bill_payments
		WHERE bill_id = ?
		ORDER BY createdAt DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []BillPayment
	for rows.Next() {
		var p BillPayment
		if err := rows.Scan(
			&p.ID,
			&p.TransactionID,
			&p.Reference,
			&p.Date,
			&p.BillID,
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

func (s *BillPaymentStore) GetByID(ctx context.Context, id int64) (*BillPayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       bill_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM bill_payments
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p BillPayment
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.TransactionID,
		&p.Reference,
		&p.Date,
		&p.BillID,
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

func (s *BillPaymentStore) GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*BillPayment, error) {
	query := `
		SELECT id, transaction_id, reference, date,
		       bill_id, user_id, account_id,
		       amount, status, createdAt, updatedAt
		FROM bill_payments
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var p BillPayment
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.TransactionID,
		&p.Reference,
		&p.Date,
		&p.BillID,
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

func (s *BillPaymentStore) Create(ctx context.Context, tx *sql.Tx, p *BillPayment) (*BillPayment, error) {
	query := `
		INSERT INTO bill_payments
		(transaction_id, reference, date, bill_id, user_id, account_id, amount, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, "1")
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		p.TransactionID,
		p.Reference,
		p.Date,
		p.BillID,
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

	p.ID = id
	return p, nil
}

func (s *BillPaymentStore) Update(ctx context.Context, tx *sql.Tx, p *BillPayment) (*BillPayment, error) {
	query := `
		UPDATE bill_payments
		SET reference = ?, date = ?, account_id = ?, amount = ?, status = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		p.Reference,
		p.Date,
		p.AccountID,
		p.Amount,
		p.Status,
		p.ID,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *BillPaymentStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bill_payments WHERE id = ?`

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

