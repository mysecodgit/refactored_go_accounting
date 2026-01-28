package store

import (
	"context"
	"database/sql"
)

type Bill struct {
	ID            int64   `json:"id"`
	BillNo        string  `json:"bill_no"`
	TransactionID int64   `json:"transaction_id"`
	BillDate      string  `json:"bill_date"`
	DueDate       string  `json:"due_date"`
	APAccountID   int64   `json:"ap_account_id"`
	UnitID        *int64  `json:"unit_id"`
	PeopleID      *int64  `json:"people_id"`
	UserID        int64   `json:"user_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	CancelReason  *string `json:"cancel_reason"`
	Status        string  `json:"status"` // enum('0','1')
	BuildingID    int64   `json:"building_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type BillStore struct {
	db *sql.DB
}

func NewBillStore(db *sql.DB) *BillStore {
	return &BillStore{db: db}
}

func (s *BillStore) GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]Bill, error) {
	query := `
		SELECT id, bill_no, transaction_id, bill_date, due_date,
		       ap_account_id, unit_id, people_id, user_id, amount,
		       description, cancel_reason, status, building_id, createdAt, updatedAt
		FROM bills
		WHERE building_id = ?
	`
	args := []any{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND DATE(bill_date) >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		query += " AND DATE(bill_date) <= ?"
		args = append(args, *endDate)
	}
	if peopleID != nil && *peopleID > 0 {
		query += " AND people_id = ?"
		args = append(args, *peopleID)
	}
	if status != nil && *status != "" {
		query += " AND status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY createdAt DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []Bill
	for rows.Next() {
		var b Bill
		if err := rows.Scan(
			&b.ID,
			&b.BillNo,
			&b.TransactionID,
			&b.BillDate,
			&b.DueDate,
			&b.APAccountID,
			&b.UnitID,
			&b.PeopleID,
			&b.UserID,
			&b.Amount,
			&b.Description,
			&b.CancelReason,
			&b.Status,
			&b.BuildingID,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bills = append(bills, b)
	}

	return bills, nil
}

func (s *BillStore) GetByID(ctx context.Context, id int64) (*Bill, error) {
	query := `
		SELECT id, bill_no, transaction_id, bill_date, due_date,
		       ap_account_id, unit_id, people_id, user_id, amount,
		       description, cancel_reason, status, building_id, createdAt, updatedAt
		FROM bills
		WHERE id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var b Bill
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.BillNo,
		&b.TransactionID,
		&b.BillDate,
		&b.DueDate,
		&b.APAccountID,
		&b.UnitID,
		&b.PeopleID,
		&b.UserID,
		&b.Amount,
		&b.Description,
		&b.CancelReason,
		&b.Status,
		&b.BuildingID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (s *BillStore) Create(ctx context.Context, tx *sql.Tx, b *Bill) (*int64, error) {
	query := `
		INSERT INTO bills
		(bill_no, transaction_id, bill_date, due_date,
		 ap_account_id, unit_id, people_id, user_id, amount,
		 description, cancel_reason, status, building_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "1", ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		b.BillNo,
		b.TransactionID,
		b.BillDate,
		b.DueDate,
		b.APAccountID,
		b.UnitID,
		b.PeopleID,
		b.UserID,
		b.Amount,
		b.Description,
		b.CancelReason,
		b.BuildingID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	b.ID = id
	return &id, nil
}

func (s *BillStore) Update(ctx context.Context, tx *sql.Tx, b *Bill) (*int64, error) {
	query := `
		UPDATE bills
		SET bill_no = ?, bill_date = ?, due_date = ?,
		    ap_account_id = ?, unit_id = ?, people_id = ?, user_id = ?,
		    amount = ?, description = ?, cancel_reason = ?, status = ?, building_id = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		b.BillNo,
		b.BillDate,
		b.DueDate,
		b.APAccountID,
		b.UnitID,
		b.PeopleID,
		b.UserID,
		b.Amount,
		b.Description,
		b.CancelReason,
		b.Status,
		b.BuildingID,
		b.ID,
	)
	if err != nil {
		return nil, err
	}

	return &b.ID, nil
}

func (s *BillStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bills WHERE id = ?`

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

