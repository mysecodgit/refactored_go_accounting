package store

import (
	"context"
	"database/sql"
)

type SalesReceipt struct {
	ID            int64  `json:"id"`
	ReceiptNo     int    `json:"receipt_no"`
	TransactionID int64  `json:"transaction_id"`
	ReceiptDate   string `json:"receipt_date"`

	UnitID   *int64 `json:"unit_id"`
	PeopleID *int64 `json:"people_id"`

	UserID    int64   `json:"user_id"`
	AccountID int64   `json:"account_id"`
	Amount    float64 `json:"amount"`

	Description  *string `json:"description"`
	CancelReason *string `json:"cancel_reason"`

	Status     int   `json:"status"` // enum('0','1')
	BuildingID int64 `json:"building_id"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// relationships
	Account Account `json:"account"`
	Unit    Unit    `json:"unit"`
	People  People  `json:"people"`
}

type SalesReceiptStore struct {
	db *sql.DB
}

func NewSalesReceiptStore(db *sql.DB) *SalesReceiptStore {
	return &SalesReceiptStore{db: db}
}

type SalesReceiptListResponse struct {
	ID            int64  `json:"id"`
	ReceiptNo     int    `json:"receipt_no"`
	TransactionID int64  `json:"transaction_id"`
	ReceiptDate   string `json:"receipt_date"`

	UnitID   *int64 `json:"unit_id"`
	PeopleID *int64 `json:"people_id"`

	UserID    int64   `json:"user_id"`
	AccountID int64   `json:"account_id"`
	Amount    float64 `json:"amount"`

	Description  *string `json:"description"`
	CancelReason *string `json:"cancel_reason"`

	Status     int   `json:"status"`
	BuildingID int64 `json:"building_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	People People `json:"people"`
	Unit   Unit   `json:"unit"`
}

func (s *SalesReceiptStore) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
	peopleID *int,
	status *string,
) ([]SalesReceiptListResponse, error) {

	query := `
		SELECT
			sr.id, sr.receipt_no, sr.transaction_id, sr.receipt_date,
			sr.unit_id, sr.people_id, sr.user_id, sr.account_id,
			sr.amount, sr.description, sr.cancel_reason,
			sr.status, sr.building_id,
			sr.createdAt, sr.updatedAt,
			p.id, p.name,
			u.id, u.name
		FROM sales_receipt sr
		LEFT JOIN people p ON p.id = sr.people_id
		LEFT JOIN units u ON u.id = sr.unit_id
		WHERE sr.building_id = ?
	`

	args := []interface{}{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND DATE(sr.receipt_date) >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil && *endDate != "" {
		query += " AND DATE(sr.receipt_date) <= ?"
		args = append(args, *endDate)
	}

	if peopleID != nil && *peopleID > 0 {
		query += " AND sr.people_id = ?"
		args = append(args, *peopleID)
	}

	if status != nil && *status != "" {
		query += " AND sr.status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY sr.createdAt DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []SalesReceiptListResponse

	for rows.Next() {
		var r SalesReceiptListResponse
		if err := rows.Scan(
			&r.ID,
			&r.ReceiptNo,
			&r.TransactionID,
			&r.ReceiptDate,
			&r.UnitID,
			&r.PeopleID,
			&r.UserID,
			&r.AccountID,
			&r.Amount,
			&r.Description,
			&r.CancelReason,
			&r.Status,
			&r.BuildingID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.People.ID,
			&r.People.Name,
			&r.Unit.ID,
			&r.Unit.Name,
		); err != nil {
			return nil, err
		}

		receipts = append(receipts, r)
	}

	return receipts, nil
}

func (s *SalesReceiptStore) GetByID(ctx context.Context, id int64) (*SalesReceipt, error) {
	query := `
		SELECT
			sr.id, sr.receipt_no, sr.transaction_id, sr.receipt_date,
			sr.unit_id, sr.people_id, sr.user_id, sr.account_id,
			sr.amount, sr.description, sr.cancel_reason,
			sr.status, sr.building_id,
			sr.createdAt, sr.updatedAt,
			a.account_name,
			u.name,
			p.name
		FROM sales_receipt sr
		LEFT JOIN accounts a ON a.id = sr.account_id
		LEFT JOIN units u ON u.id = sr.unit_id
		LEFT JOIN people p ON p.id = sr.people_id
		WHERE sr.id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var r SalesReceipt
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.ReceiptNo,
		&r.TransactionID,
		&r.ReceiptDate,
		&r.UnitID,
		&r.PeopleID,
		&r.UserID,
		&r.AccountID,
		&r.Amount,
		&r.Description,
		&r.CancelReason,
		&r.Status,
		&r.BuildingID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Account.AccountName,
		&r.Unit.Name,
		&r.People.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &r, nil
}

func (s *SalesReceiptStore) Create(ctx context.Context, tx *sql.Tx, r *SalesReceipt) (*int64, error) {
	query := `
		INSERT INTO sales_receipt
		(receipt_no, transaction_id, receipt_date,
		 unit_id, people_id, user_id, account_id,
		 amount, description, cancel_reason,
		 status, building_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '1', ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		r.ReceiptNo,
		r.TransactionID,
		r.ReceiptDate,
		r.UnitID,
		r.PeopleID,
		r.UserID,
		r.AccountID,
		r.Amount,
		r.Description,
		r.CancelReason,
		r.BuildingID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	r.ID = id
	return &id, nil
}

func (s *SalesReceiptStore) Update(ctx context.Context, tx *sql.Tx, r *SalesReceipt) (*int64, error) {
	query := `
		UPDATE sales_receipt
		SET receipt_no = ?, transaction_id = ?, receipt_date = ?,
		    unit_id = ?, people_id = ?, user_id = ?, account_id = ?,
		    amount = ?, description = ?, cancel_reason = ?,
		    building_id = ?, updatedAt = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		r.ReceiptNo,
		r.TransactionID,
		r.ReceiptDate,
		r.UnitID,
		r.PeopleID,
		r.UserID,
		r.AccountID,
		r.Amount,
		r.Description,
		r.CancelReason,
		r.BuildingID,
		r.ID,
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

	return &r.ID, nil
}

func (s *SalesReceiptStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sales_receipt WHERE id = ?`

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
