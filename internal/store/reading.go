package store

import (
	"context"
	"database/sql"
	"time"
)

// Reading represents a meter/usage reading
type Reading struct {
	ID            int64     `json:"id"`
	ItemID        int64     `json:"item_id"`
	UnitID        int64     `json:"unit_id"`
	LeaseID       *int64    `json:"lease_id"`
	ReadingMonth  *string   `json:"reading_month"`
	ReadingYear   *string   `json:"reading_year"`
	ReadingDate   time.Time `json:"reading_date"`
	PreviousValue *float64  `json:"previous_value"`
	CurrentValue  *float64  `json:"current_value"`
	UnitPrice     *float64  `json:"unit_price"`
	TotalAmount   *float64  `json:"total_amount"`
	Notes         *string   `json:"notes"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ReadingStore struct {
	db *sql.DB
}

func NewReadingStore(db *sql.DB) *ReadingStore {
	return &ReadingStore{db: db}
}

// GetAll returns all readings
func (s *ReadingStore) GetAll(ctx context.Context, buildingID int64) ([]Reading, error) {
	query := `
		SELECT id, item_id, unit_id, lease_id, reading_month, reading_year,
		       reading_date, previous_value, current_value, unit_price,
		       total_amount, notes, status, created_at, updated_at
		FROM readings where building_id = ?
		ORDER BY reading_date DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []Reading
	for rows.Next() {
		var r Reading
		if err := rows.Scan(
			&r.ID,
			&r.ItemID,
			&r.UnitID,
			&r.LeaseID,
			&r.ReadingMonth,
			&r.ReadingYear,
			&r.ReadingDate,
			&r.PreviousValue,
			&r.CurrentValue,
			&r.UnitPrice,
			&r.TotalAmount,
			&r.Notes,
			&r.Status,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		readings = append(readings, r)
	}

	return readings, nil
}

// GetByID returns a single reading by ID
func (s *ReadingStore) GetByID(ctx context.Context, id int64) (*Reading, error) {
	query := `
		SELECT id, item_id, unit_id, lease_id, reading_month, reading_year,
		       reading_date, previous_value, current_value, unit_price,
		       total_amount, notes, status, created_at, updated_at
		FROM readings
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var r Reading
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.ItemID,
		&r.UnitID,
		&r.LeaseID,
		&r.ReadingMonth,
		&r.ReadingYear,
		&r.ReadingDate,
		&r.PreviousValue,
		&r.CurrentValue,
		&r.UnitPrice,
		&r.TotalAmount,
		&r.Notes,
		&r.Status,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &r, nil
}

type ReadingByUnitResponse struct {
	ID            int64     `json:"id"`
	ItemName      string    `json:"item_name"`
	PreviousValue float64   `json:"previous_value"`
	CurrentValue  float64   `json:"current_value"`
	Consumption   float64   `json:"consumption"`
	UnitPrice     float64   `json:"unit_price"`
	TotalAmount   float64   `json:"total_amount"`
	ReadingDate   time.Time `json:"reading_date"`
	Item          Item      `json:"item"`
}

func (s *ReadingStore) GetAllByUnitID(ctx context.Context, unitID int64) ([]ReadingByUnitResponse, error) {
	query := `
	SELECT 
    r.id,
    i.id AS item_id,
    i.name AS item_name,
    IFNULL(r.previous_value, 0) AS previous_value,
    IFNULL(r.current_value, 0) AS current_value,
    IFNULL(r.current_value, 0) - IFNULL(r.previous_value, 0) AS consumption,
    r.unit_price,
    r.total_amount,
    r.reading_date
FROM readings r
LEFT JOIN items i ON r.item_id = i.id
WHERE r.unit_id = ?
  AND r.status = '1'
ORDER BY r.reading_date DESC

	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []ReadingByUnitResponse
	for rows.Next() {
		var r ReadingByUnitResponse

		if err := rows.Scan(
			&r.ID,
			&r.Item.ID,
			&r.Item.Name,
			&r.PreviousValue,
			&r.CurrentValue,
			&r.Consumption,
			&r.UnitPrice,
			&r.TotalAmount,
			&r.ReadingDate,
		); err != nil {
			return nil, err
		}

		// Optional: keep ItemName in sync
		r.ItemName = r.Item.Name

		readings = append(readings, r)
	}

	return readings, nil
}

// Create inserts a new reading
func (s *ReadingStore) Create(ctx context.Context, r *Reading) error {
	query := `
		INSERT INTO readings (
			item_id, unit_id, lease_id, reading_month, reading_year,
			reading_date, previous_value, current_value, unit_price,
			total_amount, notes, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		r.ItemID,
		r.UnitID,
		r.LeaseID,
		r.ReadingMonth,
		r.ReadingYear,
		r.ReadingDate,
		r.PreviousValue,
		r.CurrentValue,
		r.UnitPrice,
		r.TotalAmount,
		r.Notes,
		r.Status,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = id
	return nil
}

// Update modifies an existing reading
func (s *ReadingStore) Update(ctx context.Context, r *Reading) error {
	query := `
		UPDATE readings
		SET item_id = ?,
		    unit_id = ?,
		    lease_id = ?,
		    reading_month = ?,
		    reading_year = ?,
		    reading_date = ?,
		    previous_value = ?,
		    current_value = ?,
		    unit_price = ?,
		    total_amount = ?,
		    notes = ?,
		    status = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		r.ItemID,
		r.UnitID,
		r.LeaseID,
		r.ReadingMonth,
		r.ReadingYear,
		r.ReadingDate,
		r.PreviousValue,
		r.CurrentValue,
		r.UnitPrice,
		r.TotalAmount,
		r.Notes,
		r.Status,
		r.ID,
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

// Delete removes a reading by ID
func (s *ReadingStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM readings WHERE id = ?`

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
