package store

import (
	"context"
	"database/sql"
)

type Item struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"` // inventory | non inventory | service | discount
	Description    string    `json:"description"`

	AssetAccount   *int64 `json:"asset_account"`
	IncomeAccount  *int64 `json:"income_account"`
	COGSAccount    *int64 `json:"cogs_account"`
	ExpenseAccount *int64 `json:"expense_account"`

	OnHand     float64   `json:"on_hand"`
	AvgCost    float64   `json:"avg_cost"`
	Date       string `json:"date"`
	BuildingID int64     `json:"building_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ItemStore struct {
	db *sql.DB
}

func NewItemStore(db *sql.DB) *ItemStore {
	return &ItemStore{db: db}
}

func (s *ItemStore) GetAll(ctx context.Context, buildingID int64) ([]Item, error) {
	query := `
		SELECT id, name, type, description,
		       asset_account, income_account, cogs_account, expense_account,
		       on_hand, avg_cost, date, building_id,
		       created_at, updated_at
		FROM items
		WHERE building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Description,
			&i.AssetAccount,
			&i.IncomeAccount,
			&i.COGSAccount,
			&i.ExpenseAccount,
			&i.OnHand,
			&i.AvgCost,
			&i.Date,
			&i.BuildingID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (s *ItemStore) GetByID(ctx context.Context, id int64) (*Item, error) {
	query := `
		SELECT id, name, type, description,
		       asset_account, income_account, cogs_account, expense_account,
		       on_hand, avg_cost, date, building_id,
		       created_at, updated_at
		FROM items
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var i Item
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Description,
		&i.AssetAccount,
		&i.IncomeAccount,
		&i.COGSAccount,
		&i.ExpenseAccount,
		&i.OnHand,
		&i.AvgCost,
		&i.Date,
		&i.BuildingID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &i, nil
}

func (s *ItemStore) Create(ctx context.Context, i *Item) error {
	query := `
		INSERT INTO items
		(name, type, description, asset_account, income_account, cogs_account,
		 expense_account, on_hand, avg_cost, date, building_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		i.Name,
		i.Type,
		i.Description,
		i.AssetAccount,
		i.IncomeAccount,
		i.COGSAccount,
		i.ExpenseAccount,
		i.OnHand,
		i.AvgCost,
		i.Date,
		i.BuildingID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	i.ID = id
	return nil
}

func (s *ItemStore) Update(ctx context.Context, i *Item) error {
	query := `
		UPDATE items
		SET name = ?, type = ?, description = ?,
		    asset_account = ?, income_account = ?, cogs_account = ?, expense_account = ?,
		    on_hand = ?, avg_cost = ?, date = ?, building_id = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		i.Name,
		i.Type,
		i.Description,
		i.AssetAccount,
		i.IncomeAccount,
		i.COGSAccount,
		i.ExpenseAccount,
		i.OnHand,
		i.AvgCost,
		i.Date,
		i.BuildingID,
		i.ID,
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

func (s *ItemStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = ?`

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
