package store

import (
	"context"
	"database/sql"
)

type Invoice struct {
	ID            int64  `json:"id"`
	InvoiceNo     string `json:"invoice_no"`
	TransactionID int64  `json:"transaction_id"`
	SalesDate     string `json:"sales_date"`
	DueDate       string `json:"due_date"`
	ARAccountID   int    `json:"ar_account_id"`

	UnitID   *int64 `json:"unit_id"`
	PeopleID *int64 `json:"people_id"`

	UserID       int64   `json:"user_id"`
	Amount       float64 `json:"amount"`
	AmountCents  int64   `json:"amount_cents"`
	Description  string  `json:"description"`
	CancelReason *string `json:"cancel_reason"`

	Status     *int `json:"status"` // enum('0','1')
	BuildingID int64  `json:"building_id"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// relationships
	ARAccount Account `json:"ar_account"`
	Unit Unit `json:"unit"`
	People People `json:"people"`
}

type InvoiceStore struct {
	db *sql.DB
}

func NewInvoiceStore(db *sql.DB) *InvoiceStore {
	return &InvoiceStore{db: db}
}



type InvoiceSummary struct {
	ID                  int     `json:"id"`
	InvoiceNo           string  `json:"invoice_no"`
	TransactionID       int     `json:"transaction_id"`
	SalesDate           string  `json:"sales_date"`
	DueDate             string  `json:"due_date"`
	ARAccountID         *int    `json:"ar_account_id"`
	UnitID              *int    `json:"unit_id"`
	PeopleID            *int    `json:"people_id"`
	UserID              int     `json:"user_id"`
	Amount              float64 `json:"amount"`
	AmountCents         int64   `json:"amount_cents"`
	Description         string  `json:"description"`
	CancelReason        *string `json:"cancel_reason"`
	Status              int     `json:"status"`
	BuildingID          int     `json:"building_id"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	PaidAmount          float64 `json:"paid_amount"`
	PaidAmountCents     int64   `json:"paid_amount_cents"`
	AppliedCreditsTotal float64 `json:"applied_credits_total"`
	AppliedCreditsTotalCents int64 `json:"applied_credits_total_cents"`
	AppliedDiscountsTotal float64 `json:"applied_discounts_total"`
	AppliedDiscountsTotalCents int64 `json:"applied_discounts_total_cents"`
	People              People  `json:"people"`
	Unit                Unit    `json:"unit"`
}

func (s *InvoiceStore) GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]InvoiceSummary, error) {
	query := `
		SELECT 
			i.id, i.invoice_no, i.transaction_id, i.sales_date, i.due_date, 
			i.ar_account_id, i.unit_id, i.people_id, i.user_id, i.amount, i.amount_cents,
			i.description, i.cancel_reason, i.status, i.building_id, 
			i.createdAt, i.updatedAt,
			COALESCE((
				SELECT SUM(ip.amount)
				FROM invoice_payments ip 
				WHERE ip.invoice_id = i.id AND ip.status = '1'
			), 0) as paid_amount,
			COALESCE((
				SELECT SUM(ip.amount_cents)
				FROM invoice_payments ip 
				WHERE ip.invoice_id = i.id AND ip.status = '1'
			), 0) as paid_amount_cents,
			COALESCE((
				SELECT SUM(iac.amount) 
				FROM invoice_applied_credits iac 
				WHERE iac.invoice_id = i.id AND iac.status = '1'
			), 0) as applied_credits_total,
			COALESCE((
				SELECT SUM(iac.amount_cents) 
				FROM invoice_applied_credits iac 
				WHERE iac.invoice_id = i.id AND iac.status = '1'
			), 0) as applied_credits_total_cents,
			COALESCE((
				SELECT SUM(iad.amount) 
				FROM invoice_applied_discounts iad 
				WHERE iad.invoice_id = i.id AND iad.status = '1'
			), 0) as applied_discounts_total,
			COALESCE((
				SELECT SUM(iad.amount_cents) 
				FROM invoice_applied_discounts iad 
				WHERE iad.invoice_id = i.id AND iad.status = '1'
			), 0) as applied_discounts_total_cents,
			p.name as people_name,
			p.id as people_id,
			u.id as unit_id,
			u.name as unit_name
		FROM invoices i
		LEFT JOIN people p ON p.id = i.people_id
		LEFT JOIN units u ON u.id = i.unit_id
		WHERE i.building_id = ?
	`

	args := []interface{}{buildingID}

	// Add filters
	if startDate != nil && *startDate != "" {
		query += " AND DATE(i.sales_date) >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil && *endDate != "" {
		query += " AND DATE(i.sales_date) <= ?"
		args = append(args, *endDate)
	}

	if peopleID != nil && *peopleID > 0 {
		query += " AND i.people_id = ?"
		args = append(args, *peopleID)
	}

	if status != nil && *status != "" {
		query += " AND i.status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY i.createdAt DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []InvoiceSummary
	for rows.Next() {
		var invoice InvoiceSummary
		if err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNo, &invoice.TransactionID, &invoice.SalesDate, &invoice.DueDate,
			&invoice.ARAccountID, &invoice.UnitID, &invoice.PeopleID, &invoice.UserID, &invoice.Amount, &invoice.AmountCents,
			&invoice.Description, &invoice.CancelReason, &invoice.Status, &invoice.BuildingID,
			&invoice.CreatedAt, &invoice.UpdatedAt,
			&invoice.PaidAmount,&invoice.PaidAmountCents, &invoice.AppliedCreditsTotal,
			&invoice.AppliedCreditsTotalCents, &invoice.AppliedDiscountsTotal,
			&invoice.AppliedDiscountsTotalCents,
			&invoice.People.Name, &invoice.People.ID,
			&invoice.Unit.ID, &invoice.Unit.Name,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (s *InvoiceStore) GetByID(ctx context.Context, id int64) (*Invoice, error) {
	query := `
		SELECT i.id, i.invoice_no, i.transaction_id, i.sales_date, i.due_date,
		       i.ar_account_id, i.unit_id, i.people_id, i.user_id,
		       i.amount, i.amount_cents, i.description, i.cancel_reason, i.status,
		       i.building_id, i.createdAt, i.updatedAt,
			   a.account_name, u.name unit_name, p.name people_name
		FROM invoices i
		LEFT JOIN accounts a ON a.id = i.ar_account_id
		LEFT JOIN units u ON u.id = i.unit_id
		LEFT JOIN people p ON p.id = i.people_id
		WHERE i.id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var i Invoice
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID,
		&i.InvoiceNo,
		&i.TransactionID,
		&i.SalesDate,
		&i.DueDate,
		&i.ARAccountID,
		&i.UnitID,
		&i.PeopleID,
		&i.UserID,
		&i.Amount,
		&i.AmountCents,
		&i.Description,
		&i.CancelReason,
		&i.Status,
		&i.BuildingID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ARAccount.AccountName,
		&i.Unit.Name,
		&i.People.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &i, nil
}

func (s *InvoiceStore) Create(ctx context.Context, tx *sql.Tx, i *Invoice) (*int64, error) {
	query := `
		INSERT INTO invoices
		(invoice_no, transaction_id, sales_date, due_date,
		 ar_account_id, unit_id, people_id, user_id,
		 amount, amount_cents, description, cancel_reason, status, building_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "1", ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		i.InvoiceNo,
		i.TransactionID,
		i.SalesDate,
		i.DueDate,
		i.ARAccountID,
		i.UnitID,
		i.PeopleID,
		i.UserID,
		i.Amount,
		i.AmountCents,
		i.Description,
		i.CancelReason,
		i.BuildingID,
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

func (s *InvoiceStore) Update(ctx context.Context, tx *sql.Tx, i *Invoice) (*int64, error) {
	query := `
		UPDATE invoices
		SET invoice_no = ?, transaction_id = ?, sales_date = ?, due_date = ?,
		    ar_account_id = ?, unit_id = ?, people_id = ?, user_id = ?,
		    amount = ?, amount_cents = ?, description = ?, cancel_reason = ?,
		    building_id = ?, updatedAt = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		i.InvoiceNo,
		i.TransactionID,
		i.SalesDate,
		i.DueDate,
		i.ARAccountID,
		i.UnitID,
		i.PeopleID,
		i.UserID,
		i.Amount,
		i.AmountCents,
		i.Description,
		i.CancelReason,
		i.BuildingID,
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

func (s *InvoiceStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM invoices WHERE id = ?`

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
