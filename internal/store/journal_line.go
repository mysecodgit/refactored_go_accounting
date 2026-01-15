package store

import (
	"context"
	"database/sql"
)

type JournalLine struct {
	ID          int64    `json:"id"`
	JournalID   int64    `json:"journal_id"`
	AccountID   int64    `json:"account_id"`
	UnitID      *int64   `json:"unit_id,omitempty"`
	PeopleID    *int64   `json:"people_id,omitempty"`
	Description *string  `json:"description,omitempty"`
	Debit       float64  `json:"debit,omitempty"`
	Credit      float64  `json:"credit,omitempty"`
}

type JournalLineStore struct {
	db *sql.DB
}

func NewJournalLineStore(db *sql.DB) *JournalLineStore {
	return &JournalLineStore{db: db}
}

func (s *JournalLineStore) GetAllByJournalID(ctx context.Context, journalID int64) ([]JournalLine, error) {
	query := `SELECT id, journal_id, account_id, unit_id, people_id, description, debit, credit
			  FROM journal_lines
			  WHERE journal_id = ?
			  ORDER BY id ASC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, journalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []JournalLine
	for rows.Next() {
		var l JournalLine
		if err := rows.Scan(
			&l.ID,
			&l.JournalID,
			&l.AccountID,
			&l.UnitID,
			&l.PeopleID,
			&l.Description,
			&l.Debit,
			&l.Credit,
		); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}

	return lines, nil
}

func (s *JournalLineStore) GetByID(ctx context.Context, id int64) (*JournalLine, error) {
	query := `SELECT id, journal_id, account_id, unit_id, people_id, description, debit, credit
			  FROM journal_lines
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var l JournalLine
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.JournalID,
		&l.AccountID,
		&l.UnitID,
		&l.PeopleID,
		&l.Description,
		&l.Debit,
		&l.Credit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &l, nil
}

func (s *JournalLineStore) GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*JournalLine, error) {
	query := `SELECT id, journal_id, account_id, unit_id, people_id, description, debit, credit
			  FROM journal_lines
			  WHERE id = ?`

	var l JournalLine
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.JournalID,
		&l.AccountID,
		&l.UnitID,
		&l.PeopleID,
		&l.Description,
		&l.Debit,
		&l.Credit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (s *JournalLineStore) Create(ctx context.Context, tx *sql.Tx, l *JournalLine) (*JournalLine, error) {
	query := `INSERT INTO journal_lines
			  (journal_id, account_id, unit_id, people_id, description, debit, credit)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		l.JournalID,
		l.AccountID,
		l.UnitID,
		l.PeopleID,
		l.Description,
		l.Debit,
		l.Credit,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	created, err := s.GetByIDTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *JournalLineStore) Update(ctx context.Context, tx *sql.Tx, l *JournalLine) (*JournalLine, error) {
	query := `UPDATE journal_lines
			  SET journal_id = ?, account_id = ?, unit_id = ?, people_id = ?, description = ?, debit = ?, credit = ?
			  WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query,
		l.JournalID,
		l.AccountID,
		l.UnitID,
		l.PeopleID,
		l.Description,
		l.Debit,
		l.Credit,
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

	updated, err := s.GetByIDTx(ctx, tx, l.ID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *JournalLineStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM journal_lines WHERE id = ?`

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

func (s *JournalLineStore) DeleteByJournalID(ctx context.Context, tx *sql.Tx, journalID int64) error {
	query := `DELETE FROM journal_lines WHERE journal_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, query, journalID)
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
