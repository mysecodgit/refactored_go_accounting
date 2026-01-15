package store

import (
	"context"
	"database/sql"
)

type Journal struct {
	ID            int64      `json:"id"`
	TransactionID int64      `json:"transaction_id"`
	Reference     string     `json:"reference"`
	JournalDate   string  `json:"journal_date"`
	BuildingID    int64      `json:"building_id"`
	Memo          *string    `json:"memo,omitempty"`
	TotalAmount   *float64   `json:"total_amount,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type JournalStore struct {
	db *sql.DB
}

func NewJournalStore(db *sql.DB) *JournalStore {
	return &JournalStore{db: db}
}

func (s *JournalStore) GetAll(ctx context.Context, buildingID int64, startDate, endDate *string) ([]Journal, error) {
	query := `SELECT id, transaction_id, reference, journal_date, building_id, memo, total_amount, created_at
			  FROM journal
			  WHERE building_id = ?`


			  args := []interface{}{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND journal_date >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		query += " AND journal_date <= ?"
		args = append(args, *endDate)
	}

	query += " ORDER BY journal_date DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var j Journal
		if err := rows.Scan(
			&j.ID,
			&j.TransactionID,
			&j.Reference,
			&j.JournalDate,
			&j.BuildingID,
			&j.Memo,
			&j.TotalAmount,
			&j.CreatedAt,
		); err != nil {
			return nil, err
		}
		journals = append(journals, j)
	}

	return journals, nil
}

func (s *JournalStore) GetByID(ctx context.Context, id int64) (*Journal, error) {
	query := `SELECT id, transaction_id, reference, journal_date, building_id, memo, total_amount, created_at
			  FROM journal
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var j Journal
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID,
		&j.TransactionID,
		&j.Reference,
		&j.JournalDate,
		&j.BuildingID,
		&j.Memo,
		&j.TotalAmount,
		&j.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &j, nil
}

func (s *JournalStore) GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*Journal, error) {
	query := `SELECT id, transaction_id, reference, journal_date, building_id, memo, total_amount, created_at
			  FROM journal
			  WHERE id = ?`

	var j Journal
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&j.ID,
		&j.TransactionID,
		&j.Reference,
		&j.JournalDate,
		&j.BuildingID,
		&j.Memo,
		&j.TotalAmount,
		&j.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &j, nil
}

func (s *JournalStore) GetByTransactionID(ctx context.Context, tx *sql.Tx, transactionID int64) (*Journal, error) {
	query := `SELECT id, transaction_id, reference, journal_date, building_id, memo, total_amount, created_at
			  FROM journal
			  WHERE transaction_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var j Journal
	err := tx.QueryRowContext(ctx, query, transactionID).Scan(
		&j.ID,
		&j.TransactionID,
		&j.Reference,
		&j.JournalDate,
		&j.BuildingID,
		&j.Memo,
		&j.TotalAmount,
		&j.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &j, nil
}

func (s *JournalStore) Create(ctx context.Context, tx *sql.Tx, j *Journal) (*Journal, error) {
	query := `INSERT INTO journal
			  (transaction_id, reference, journal_date, building_id, memo, total_amount)
			  VALUES (?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		j.TransactionID,
		j.Reference,
		j.JournalDate,
		j.BuildingID,
		j.Memo,
		j.TotalAmount,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	// Fetch the fully populated record (includes created_at)
	created, err := s.GetByIDTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *JournalStore) Update(ctx context.Context, tx *sql.Tx, j *Journal) (*Journal, error) {
	query := `UPDATE journal
			  SET transaction_id = ?, reference = ?, journal_date = ?, building_id = ?, memo = ?, total_amount = ?
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		j.TransactionID,
		j.Reference,
		j.JournalDate,
		j.BuildingID,
		j.Memo,
		j.TotalAmount,
		j.ID,
	)
	if err != nil {
		return nil, err
	}

	

	updated, err := s.GetByIDTx(ctx, tx, j.ID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *JournalStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM journal WHERE id = ?`

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
