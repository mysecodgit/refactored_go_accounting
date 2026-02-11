package store

import (
	"context"
	"database/sql"
)

type Check struct {
	ID               int64   `json:"id"`
	TransactionID    int64   `json:"transaction_id"`
	CheckDate        string  `json:"check_date"`
	ReferenceNumber  string  `json:"reference_number"`
	PaymentAccountID int64   `json:"payment_account_id"`
	BuildingID       int64   `json:"building_id"`
	Memo             *string `json:"memo"`
	TotalAmount      float64 `json:"total_amount"`
	AmountCents      int64   `json:"amount_cents"`
	CreatedAt        string  `json:"created_at"`
}

type CheckStore struct {
	db *sql.DB
}

func NewCheckStore(db *sql.DB) *CheckStore {
	return &CheckStore{db: db}
}

func (s *CheckStore) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
) ([]Check, error) {
	query := `SELECT id, transaction_id, check_date, reference_number,
                     payment_account_id, building_id, memo, total_amount, amount_cents, created_at
              FROM checks
              WHERE building_id = ?`
	args := []interface{}{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND check_date >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		query += " AND check_date <= ?"
		args = append(args, *endDate)
	}

	query += " ORDER BY created_at DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []Check
	for rows.Next() {
		var c Check
		if err := rows.Scan(
			&c.ID,
			&c.TransactionID,
			&c.CheckDate,
			&c.ReferenceNumber,
			&c.PaymentAccountID,
			&c.BuildingID,
			&c.Memo,
			&c.TotalAmount,
			&c.AmountCents,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		checks = append(checks, c)
	}

	return checks, nil
}

func (s *CheckStore) GetByID(ctx context.Context, id int64) (*Check, error) {
	query := `SELECT id, transaction_id, check_date, reference_number,
                     payment_account_id, building_id, memo, total_amount, amount_cents, created_at
              FROM checks
              WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var c Check
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.TransactionID,
		&c.CheckDate,
		&c.ReferenceNumber,
		&c.PaymentAccountID,
		&c.BuildingID,
		&c.Memo,
		&c.TotalAmount,
		&c.AmountCents,
		&c.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (s *CheckStore) Create(ctx context.Context, tx *sql.Tx, c *Check) (*int64, error) {
	query := `INSERT INTO checks
              (transaction_id, check_date, reference_number, payment_account_id, building_id, memo, total_amount, amount_cents)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		c.TransactionID,
		c.CheckDate,
		c.ReferenceNumber,
		c.PaymentAccountID,
		c.BuildingID,
		c.Memo,
		c.TotalAmount,
		c.AmountCents,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	c.ID = id
	return &id, nil
}

func (s *CheckStore) Update(ctx context.Context, tx *sql.Tx, c *Check) (*int64, error) {
	query := `UPDATE checks
              SET transaction_id = ?, check_date = ?, reference_number = ?,
                  payment_account_id = ?, building_id = ?, memo = ?, total_amount = ?, amount_cents = ?
              WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		c.TransactionID,
		c.CheckDate,
		c.ReferenceNumber,
		c.PaymentAccountID,
		c.BuildingID,
		c.Memo,
		c.TotalAmount,
		c.AmountCents,
		c.ID,
	)
	if err != nil {
		return nil, err
	}

	// rows, err := result.RowsAffected()
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("+-++++++++++++++++++++++++++++++++++++++++++++++ Rows", rows)
	// if rows == 0 {
	// 	return nil, ErrNotFound
	// }

	return &c.ID, nil
}

func (s *CheckStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM checks WHERE id = ?`

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
