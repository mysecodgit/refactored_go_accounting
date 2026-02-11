package store

import (
	"context"
	"database/sql"
	"fmt"
)

type ExpenseLine struct {
	ID          int64   `json:"id"`
	CheckID     int64   `json:"check_id"`
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id,omitempty"`
	PeopleID    *int64  `json:"people_id,omitempty"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	AmountCents int64   `json:"amount_cents"`
}

type ExpenseLineStore struct {
	db *sql.DB
}

func NewExpenseLineStore(db *sql.DB) *ExpenseLineStore {
	return &ExpenseLineStore{db: db}
}

func (s *ExpenseLineStore) GetAllByCheckID(ctx context.Context, checkID int64) ([]ExpenseLine, error) {
	query := `SELECT id, check_id, account_id, unit_id, people_id, description, amount, amount_cents
			  FROM expense_lines
			  WHERE check_id = ?
			  ORDER BY id ASC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, checkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []ExpenseLine
	for rows.Next() {
		var l ExpenseLine
		if err := rows.Scan(
			&l.ID,
			&l.CheckID,
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

func (s *ExpenseLineStore) GetByID(ctx context.Context, id int64) (*ExpenseLine, error) {
	query := `SELECT id, check_id, account_id, unit_id, people_id, description, amount, amount_cents
			  FROM expense_lines
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var l ExpenseLine
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.CheckID,
		&l.AccountID,
		&l.UnitID,
		&l.PeopleID,
		&l.Description,
		&l.Amount,
		&l.AmountCents,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &l, nil
}

func (s *ExpenseLineStore) Create(ctx context.Context, tx *sql.Tx, l *ExpenseLine) (*int64, error) {
	query := `INSERT INTO expense_lines
			  (check_id, account_id, unit_id, people_id, description, amount, amount_cents)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		l.CheckID,
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

func (s *ExpenseLineStore) Update(ctx context.Context, tx *sql.Tx, l *ExpenseLine) (*int64, error) {
	query := `UPDATE expense_lines
			  SET check_id = ?, account_id = ?, unit_id = ?, people_id = ?, description = ?, amount = ?, amount_cents = ?
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		l.CheckID,
		l.AccountID,
		l.UnitID,
		l.PeopleID,
		l.Description,
		l.Amount,
		l.AmountCents,
		l.ID,
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

	return &l.ID, nil
}

func (s *ExpenseLineStore) Delete(ctx context.Context, id int64) error {
	query := `delete from expense_lines WHERE id = ?`

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

func (s *ExpenseLineStore) DeleteByCheckID(ctx context.Context, tx *sql.Tx, checkID int64) error {
	fmt.Println("+-++++++++++++++++++++++++++++++++++++++++++++++ DeleteByCheckID", checkID)
	query := `delete from expense_lines WHERE check_id = ?`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	result, err := tx.ExecContext(ctx, query, checkID)
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
