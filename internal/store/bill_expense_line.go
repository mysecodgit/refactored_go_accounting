package store

import (
	"context"
	"database/sql"
)

type BillExpenseLine struct {
	ID          int64   `json:"id"`
	BillID      int64   `json:"bill_id"`
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id"`
	PeopleID    *int64  `json:"people_id"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	AmountCents int64   `json:"amount_cents"`
}

type BillExpenseLineStore struct {
	db *sql.DB
}

func NewBillExpenseLineStore(db *sql.DB) *BillExpenseLineStore {
	return &BillExpenseLineStore{db: db}
}

func (s *BillExpenseLineStore) GetAllByBillID(ctx context.Context, billID int64) ([]BillExpenseLine, error) {
	query := `
		SELECT id, bill_id, account_id, unit_id, people_id, description, amount, amount_cents
		FROM bill_expense_lines
		WHERE bill_id = ?
		ORDER BY id ASC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []BillExpenseLine
	for rows.Next() {
		var l BillExpenseLine
		if err := rows.Scan(
			&l.ID,
			&l.BillID,
			&l.AccountID,
			&l.UnitID,
			&l.PeopleID,
			&l.Description,
			&l.Amount,
			&l.AmountCents,
		); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, nil
}

func (s *BillExpenseLineStore) Create(ctx context.Context, tx *sql.Tx, l *BillExpenseLine) (*int64, error) {
	query := `
		INSERT INTO bill_expense_lines
		(bill_id, account_id, unit_id, people_id, description, amount, amount_cents)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		l.BillID,
		l.AccountID,
		l.UnitID,
		l.PeopleID,
		l.Description,
		l.Amount,
		l.AmountCents,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	l.ID = id
	return &id, nil
}

func (s *BillExpenseLineStore) DeleteByBillID(ctx context.Context, tx *sql.Tx, billID int64) error {
	query := `DELETE FROM bill_expense_lines WHERE bill_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, billID)
	return err
}

