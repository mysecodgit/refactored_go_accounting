package store

import (
	"context"
	"database/sql"
	"strconv"
)

type Lease struct {
	ID int64 `json:"id"`

	PeopleID   int64 `json:"people_id"`
	BuildingID int64 `json:"building_id"`
	UnitID     int64 `json:"unit_id"`

	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date"`

	RentAmount    float64 `json:"rent_amount"`
	DepositAmount float64 `json:"deposit_amount"`
	ServiceAmount float64 `json:"service_amount"`

	LeaseTerms string `json:"lease_terms"`
	Status     int    `json:"status"` // enum('0','1')

	// relationships
	People People `json:"people"`
	Unit   Unit   `json:"unit"`
}

type LeaseStore struct {
	db *sql.DB
}

func NewLeaseStore(db *sql.DB) *LeaseStore {
	return &LeaseStore{db: db}
}

func (s *LeaseStore) GetAll(
	ctx context.Context,
	buildingID int64,
	peopleID *int64,
	unitID *int64,
	status *string,
) ([]Lease, error) {

	query := `
		SELECT
			l.id, l.people_id, l.building_id, l.unit_id,
			l.start_date, l.end_date,
			l.rent_amount, l.deposit_amount, l.service_amount,
			l.lease_terms, l.status,
			p.id, p.name,
			u.id, u.name
		FROM leases l
		LEFT JOIN people p ON p.id = l.people_id
		LEFT JOIN units u ON u.id = l.unit_id
		WHERE l.building_id = ?
	`

	args := []any{buildingID}

	if peopleID != nil && *peopleID > 0 {
		query += " AND l.people_id = ?"
		args = append(args, *peopleID)
	}

	if unitID != nil && *unitID > 0 {
		query += " AND l.unit_id = ?"
		args = append(args, *unitID)
	}

	if status != nil && *status != "" {
		query += " AND l.status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY l.start_date DESC"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []Lease

	for rows.Next() {
		var l Lease
		if err := rows.Scan(
			&l.ID,
			&l.PeopleID,
			&l.BuildingID,
			&l.UnitID,
			&l.StartDate,
			&l.EndDate,
			&l.RentAmount,
			&l.DepositAmount,
			&l.ServiceAmount,
			&l.LeaseTerms,
			&l.Status,
			&l.People.ID,
			&l.People.Name,
			&l.Unit.ID,
			&l.Unit.Name,
		); err != nil {
			return nil, err
		}

		leases = append(leases, l)
	}

	return leases, nil
}

func (s *LeaseStore) GetByID(ctx context.Context, id int64) (*Lease, error) {
	query := `
		SELECT
			l.id, l.people_id, l.building_id, l.unit_id,
			l.start_date, l.end_date,
			l.rent_amount, l.deposit_amount, l.service_amount,
			l.lease_terms, l.status,
			p.name,
			u.name
		FROM leases l
		LEFT JOIN people p ON p.id = l.people_id
		LEFT JOIN units u ON u.id = l.unit_id
		WHERE l.id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var l Lease
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.PeopleID,
		&l.BuildingID,
		&l.UnitID,
		&l.StartDate,
		&l.EndDate,
		&l.RentAmount,
		&l.DepositAmount,
		&l.ServiceAmount,
		&l.LeaseTerms,
		&l.Status,
		&l.People.Name,
		&l.Unit.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &l, nil
}

func (s *LeaseStore) Create(ctx context.Context, tx *sql.Tx, l *Lease) (*int64, error) {
	query := `
		INSERT INTO leases
		(people_id, building_id, unit_id,
		 start_date, end_date,
		 rent_amount, deposit_amount, service_amount,
		 lease_terms, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		l.PeopleID,
		l.BuildingID,
		l.UnitID,
		l.StartDate,
		l.EndDate,
		l.RentAmount,
		l.DepositAmount,
		l.ServiceAmount,
		l.LeaseTerms,
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

func (s *LeaseStore) Update(ctx context.Context, tx *sql.Tx, l *Lease) (*int64, error) {
	query := `
		UPDATE leases
		SET people_id = ?, building_id = ?, unit_id = ?,
		    start_date = ?, end_date = ?,
		    rent_amount = ?, deposit_amount = ?, service_amount = ?,
		    lease_terms = ?, status = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		l.PeopleID,
		l.BuildingID,
		l.UnitID,
		l.StartDate,
		l.EndDate,
		l.RentAmount,
		l.DepositAmount,
		l.ServiceAmount,
		l.LeaseTerms,
		strconv.Itoa(l.Status),
		l.ID,
	)
	if err != nil {
		return nil, err
	}

	return &l.ID, nil
}

func (s *LeaseStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM leases WHERE id = ?`

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
