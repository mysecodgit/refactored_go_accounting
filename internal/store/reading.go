package store

import (
	"context"
	"database/sql"
)

// Reading represents a meter/usage reading
type Reading struct {
	ID            int64    `json:"id"`
	ItemID        int64    `json:"item_id"`
	UnitID        int64    `json:"unit_id"`
	LeaseID       *int64   `json:"lease_id"`
	ReadingMonth  *string  `json:"reading_month"`
	ReadingYear   *string  `json:"reading_year"`
	ReadingDate   string   `json:"reading_date"`
	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
	UnitPrice     *float64 `json:"unit_price"`
	TotalAmount   *float64 `json:"total_amount"`
	Notes         *string  `json:"notes"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`

	// relationships
	Item       Item    `json:"item"`
	Unit       Unit    `json:"unit"`
	PeopleName *string `json:"people_name"`
}

type ReadingStore struct {
	db *sql.DB
}

func NewReadingStore(db *sql.DB) *ReadingStore {
	return &ReadingStore{db: db}
}

// GetAll returns all readings
func (s *ReadingStore) GetAll(ctx context.Context, buildingID int64, status *string) ([]Reading, error) {
	query := `
		SELECT r.id, i.id as item_id, u.id as unit_id, l.id as lease_id, r.reading_month, r.reading_year,
		       r.reading_date, r.previous_value, r.current_value, r.unit_price,
		       r.total_amount, r.notes, r.status, r.created_at, r.updated_at,
			   i.name as item_name, u.name as unit_name, p.name as people_name
		FROM readings r 
		LEFT JOIN items i ON r.item_id = i.id
		LEFT JOIN units u ON r.unit_id = u.id
		LEFT JOIN buildings b ON u.building_id = b.id
		LEFT JOIN leases l ON r.lease_id = l.id
		LEFT JOIN people p ON l.people_id = p.id
		where b.id = ?
		
	`

	if status != nil {
		query += " and r.status = ?"
	}

	query += " ORDER BY r.reading_date DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID, status)
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
			&r.Item.Name,
			&r.Unit.Name,
			&r.PeopleName,
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
	ID            int64   `json:"id"`
	ItemName      string  `json:"item_name"`
	PreviousValue float64 `json:"previous_value"`
	CurrentValue  float64 `json:"current_value"`
	Consumption   float64 `json:"consumption"`
	UnitPrice     float64 `json:"unit_price"`
	TotalAmount   float64 `json:"total_amount"`
	ReadingDate   string  `json:"reading_date"`
	Item          Item    `json:"item"`
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

// get latest reading by item id and unit id
func (s *ReadingStore) GetLatest(ctx context.Context, itemID int64, unitID int64) (*Reading, error) {
	query := `
		SELECT * FROM readings WHERE item_id = ? AND unit_id = ? AND status = '1' ORDER BY reading_date DESC LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, itemID, unitID)
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

	return &readings[0], nil
}

// get latest reading by item id and unit id
func (s *ReadingStore) GetLatestTx(ctx context.Context, tx *sql.Tx, itemID int64, unitID int64) (*Reading, error) {
	query := `
		SELECT * FROM readings WHERE item_id = ? AND unit_id = ? AND status = '1' ORDER BY reading_date DESC LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	row := tx.QueryRowContext(ctx, query, itemID, unitID)

	var r Reading
	err := row.Scan(
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

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Create inserts a new reading
func (s *ReadingStore) Create(ctx context.Context, tx *sql.Tx, r *Reading) error {
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

	result, err := tx.ExecContext(
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

	_, err := s.db.ExecContext(
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
