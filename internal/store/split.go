package store

import (
	"context"
	"database/sql"
	"fmt"
	
)

// Split represents a breakdown of a transaction
type Split struct {
	ID            int64    `json:"id"`
	TransactionID int64    `json:"transaction_id"`
	AccountID     int64    `json:"account_id"`
	Debit         *float64 `json:"debit"`
	Credit        *float64 `json:"credit"`
	UnitID        *int64   `json:"unit_id"`
	PeopleID      *int64   `json:"people_id"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	DebitCents    *int64   `json:"debit_cents"`
	CreditCents   *int64   `json:"credit_cents"`

	// relationships
	Account Account `json:"account"`
	Unit    Unit    `json:"unit"`
	People  People  `json:"people"`
}



type SplitStore struct {
	db *sql.DB
}

func NewSplitStore(db *sql.DB) *SplitStore {
	return &SplitStore{db: db}
}

// GetAll returns all splits for a transaction
func (s *SplitStore) GetAll(ctx context.Context, transactionID int64) ([]Split, error) {
	query := `
		SELECT sp.id, sp.transaction_id, sp.account_id, sp.debit,sp.credit,sp.unit_id,sp.people_id, sp.status, sp.created_at, sp.updated_at,
		a.account_name,u.name unit_name,p.name people_name, sp.debit_cents, sp.credit_cents
		FROM splits sp
		LEFT JOIN accounts a ON a.id = sp.account_id
		LEFT JOIN units u ON u.id = sp.unit_id
		LEFT JOIN people p ON p.id = sp.people_id
		WHERE transaction_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []Split
	for rows.Next() {
		var sp Split

		if err := rows.Scan(
			&sp.ID,
			&sp.TransactionID,
			&sp.AccountID,
			&sp.Debit,
			&sp.Credit,
			&sp.UnitID,
			&sp.PeopleID,
			&sp.Status,
			&sp.CreatedAt,
			&sp.UpdatedAt,
			&sp.Account.AccountName,
			&sp.Unit.Name,
			&sp.People.Name,
			&sp.DebitCents,
			&sp.CreditCents,
		); err != nil {
			return nil, err
		}
		splits = append(splits, sp)
	}

	return splits, nil
}

// GetByID returns a single split by ID
func (s *SplitStore) GetByID(ctx context.Context, id int64) (*Split, error) {
	query := `
		SELECT id, transaction_id, account_id, debit,credit,unit_id,people_id, status, created_at, updated_at, debit_cents, credit_cents
		FROM splits
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var sp Split
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sp.ID,
		&sp.TransactionID,
		&sp.AccountID,
		&sp.Debit,
		&sp.Credit,
		&sp.UnitID,
		&sp.PeopleID,
		&sp.Status,
		&sp.CreatedAt,
		&sp.UpdatedAt,
		&sp.DebitCents,
		&sp.CreditCents,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &sp, nil
}

func (s *SplitStore) GetByTransactionID(ctx context.Context, transactionID int64) ([]Split, error) {

	query := `
		SELECT id, transaction_id, account_id, debit,credit,unit_id,people_id, status, created_at, updated_at, debit_cents, credit_cents
		FROM splits
		WHERE transaction_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println("transactionID   ----------------- transactionID : ", transactionID)
	fmt.Println("rows   -----------------  : ", rows)

	var splits []Split
	for rows.Next() {
		var sp Split
		if err := rows.Scan(
			&sp.ID,
			&sp.TransactionID,
			&sp.AccountID,
			&sp.Debit,
			&sp.Credit,
			&sp.UnitID,
			&sp.PeopleID,
			&sp.Status,
			&sp.CreatedAt,
			&sp.UpdatedAt,
			&sp.DebitCents,
			&sp.CreditCents,
		); err != nil {
			return nil, err
		}
		splits = append(splits, sp)
	}

	return splits, nil
}

// Create inserts a new split
func (s *SplitStore) Create(ctx context.Context, tx *sql.Tx, sp *Split) error {
	query := `
		INSERT INTO splits (transaction_id, account_id, debit, credit, unit_id, people_id, status, debit_cents, credit_cents)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	fmt.Println("sp   ----------------- Status : ", sp.Status)
	fmt.Printf("Status='%s', len=%d\n", sp.Status, len(sp.Status))

	result, err := tx.ExecContext(
		ctx,
		query,
		sp.TransactionID,
		sp.AccountID,
		sp.Debit,
		sp.Credit,
		sp.UnitID,
		sp.PeopleID,
		sp.Status,
		sp.DebitCents,
		sp.CreditCents,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	sp.ID = id
	return nil
}

// Update modifies an existing split
func (s *SplitStore) Update(ctx context.Context, sp *Split) error {
	query := `
		UPDATE splits
		SET transaction_id = ?, account_id = ?, amount = ?, description = ?, updated_at = CURRENT_TIMESTAMP, debit_cents = ?, credit_cents = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		sp.TransactionID,
		sp.AccountID,
		sp.ID,
		sp.DebitCents,
		sp.CreditCents,
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

// Delete removes a split by ID
func (s *SplitStore) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `UPDATE splits SET status = '0' WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, id)
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

// Delete removes a split by ID
func (s *SplitStore) DeleteByTransactionID(ctx context.Context, tx *sql.Tx, transactionID int64) error {
	query := `UPDATE splits SET status = '0' WHERE transaction_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, transactionID)
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
