package store

import (
	"context"
	"database/sql"
)

type CreditMemo struct {
	ID               int64   `json:"id"`
	TransactionID    int64   `json:"transaction_id"`
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	UserID           int64   `json:"user_id"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int64   `json:"people_id"`
	BuildingID       int64   `json:"building_id"`
	UnitID           int64   `json:"unit_id"`
	Amount           float64 `json:"amount"`
	AmountCents      int64   `json:"amount_cents"`
	Description      string  `json:"description"`
	Status           int     `json:"status"` // enum('0','1')
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`

	// relationships
	People People `json:"people"`
	Unit   Unit   `json:"unit"`
}

type CreditMemoStore struct {
	db *sql.DB
}

func NewCreditMemoStore(db *sql.DB) *CreditMemoStore {
	return &CreditMemoStore{db: db}
}

type CreditMemoSummary struct {
	ID               int     `json:"id"`
	TransactionID    int     `json:"transaction_id"`
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	UserID           int     `json:"user_id"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int     `json:"people_id"`
	BuildingID       int     `json:"building_id"`
	UnitID           int     `json:"unit_id"`
	Amount           float64 `json:"amount"`
	AmountCents      int64   `json:"amount_cents"`
	Description      string  `json:"description"`
	Status           int     `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	People           People  `json:"people"`
	Unit             Unit    `json:"unit"`
	UsedCredits      float64 `json:"used_credits"`
	Balance          float64 `json:"balance"`
}

func (s *CreditMemoStore) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
	peopleID *int,
	status *string,
) ([]CreditMemoSummary, error) {

	query := `
		SELECT 
			cm.id, cm.transaction_id, cm.reference, cm.date,
			cm.user_id, cm.deposit_to, cm.liability_account,
			cm.people_id, cm.building_id, cm.unit_id,
			cm.amount, cm.amount_cents, cm.description, cm.status,
			cm.created_at, cm.updated_at,
			p.id, p.name,
			u.id, u.name,
			IFNULL(sum(ic.amount), 0) as used_credits,
			(cm.amount - IFNULL(sum(ic.amount), 0)) as balance
		FROM credit_memo cm
		LEFT JOIN people p ON p.id = cm.people_id
		LEFT JOIN units u ON u.id = cm.unit_id
		LEFT JOIN invoice_applied_credits ic ON ic.credit_memo_id = cm.id and ic.status = '1'
		WHERE cm.building_id = ? and cm.status = '1'
		GROUP BY p.id,cm.id
	`

	args := []interface{}{buildingID}

	if startDate != nil && *startDate != "" {
		query += " AND DATE(cm.date) >= ?"
		args = append(args, *startDate)
	}

	if endDate != nil && *endDate != "" {
		query += " AND DATE(cm.date) <= ?"
		args = append(args, *endDate)
	}

	if peopleID != nil && *peopleID > 0 {
		query += " AND cm.people_id = ?"
		args = append(args, *peopleID)
	}

	if status != nil && *status != "" {
		query += " AND cm.status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY cm.created_at DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creditMemos []CreditMemoSummary
	for rows.Next() {
		var cm CreditMemoSummary
		if err := rows.Scan(
			&cm.ID,
			&cm.TransactionID,
			&cm.Reference,
			&cm.Date,
			&cm.UserID,
			&cm.DepositTo,
			&cm.LiabilityAccount,
			&cm.PeopleID,
			&cm.BuildingID,
			&cm.UnitID,
			&cm.Amount,
			&cm.AmountCents,
			&cm.Description,
			&cm.Status,
			&cm.CreatedAt,
			&cm.UpdatedAt,
			&cm.People.ID,
			&cm.People.Name,
			&cm.Unit.ID,
			&cm.Unit.Name,
			&cm.UsedCredits,
			&cm.Balance,
		); err != nil {
			return nil, err
		}
		creditMemos = append(creditMemos, cm)
	}

	return creditMemos, nil
}

func (s *CreditMemoStore) GetByID(ctx context.Context, id int64) (*CreditMemo, error) {
	query := `
		SELECT 
			cm.id, cm.transaction_id, cm.reference, cm.date,
			cm.user_id, cm.deposit_to, cm.liability_account,
			cm.people_id, cm.building_id, cm.unit_id,
			cm.amount,cm.amount_cents, cm.description, cm.status,
			cm.created_at, cm.updated_at,
			p.name, u.name
		FROM credit_memo cm
		LEFT JOIN people p ON p.id = cm.people_id
		LEFT JOIN units u ON u.id = cm.unit_id
		WHERE cm.id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var cm CreditMemo
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&cm.ID,
		&cm.TransactionID,
		&cm.Reference,
		&cm.Date,
		&cm.UserID,
		&cm.DepositTo,
		&cm.LiabilityAccount,
		&cm.PeopleID,
		&cm.BuildingID,
		&cm.UnitID,
		&cm.Amount,
		&cm.AmountCents,
		&cm.Description,
		&cm.Status,
		&cm.CreatedAt,
		&cm.UpdatedAt,
		&cm.People.Name,
		&cm.Unit.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &cm, nil
}

func (s *CreditMemoStore) GetByPeopleID(ctx context.Context, peopleID int64) ([]CreditMemo, error) {
	query := `
		SELECT * FROM credit_memo WHERE people_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, peopleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []CreditMemo
	for rows.Next() {
		var cm CreditMemo
		if err := rows.Scan(
			&cm.ID,
			&cm.TransactionID,
			&cm.Reference,
			&cm.Date,
			&cm.UserID,
			&cm.DepositTo,
			&cm.LiabilityAccount,
			&cm.PeopleID,
			&cm.BuildingID,
			&cm.UnitID,
			&cm.Amount,
			&cm.Description,
			&cm.Status,
			&cm.CreatedAt,
			&cm.UpdatedAt,
			&cm.AmountCents, // import to be last cause i used select * 
		); err != nil {
			return nil, err
		}
		credits = append(credits, cm)
	}

	return credits, nil
}

func (s *CreditMemoStore) Create(ctx context.Context, tx *sql.Tx, cm *CreditMemo) (*int64, error) {
	query := `
		INSERT INTO credit_memo
		(transaction_id, reference, date, user_id,
		 deposit_to, liability_account, people_id,
		 building_id, unit_id, amount,amount_cents, description, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, '1')
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		cm.TransactionID,
		cm.Reference,
		cm.Date,
		cm.UserID,
		cm.DepositTo,
		cm.LiabilityAccount,
		cm.PeopleID,
		cm.BuildingID,
		cm.UnitID,
		cm.Amount,
		cm.AmountCents,
		cm.Description,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	cm.ID = id
	return &id, nil
}

func (s *CreditMemoStore) Update(ctx context.Context, tx *sql.Tx, cm *CreditMemo) (*int64, error) {
	query := `
		UPDATE credit_memo
		SET transaction_id = ?, reference = ?, date = ?, user_id = ?,
		    deposit_to = ?, liability_account = ?, people_id = ?,
		    building_id = ?, unit_id = ?, amount = ?,amount_cents = ?, description = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		cm.TransactionID,
		cm.Reference,
		cm.Date,
		cm.UserID,
		cm.DepositTo,
		cm.LiabilityAccount,
		cm.PeopleID,
		cm.BuildingID,
		cm.UnitID,
		cm.Amount,
		cm.AmountCents,
		cm.Description,
		cm.ID,
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

	return &cm.ID, nil
}

func (s *CreditMemoStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM credit_memo WHERE id = ?`

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
