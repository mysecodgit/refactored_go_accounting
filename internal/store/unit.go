package store

import (
	"context"
	"database/sql"
	"time"
)

type Unit struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	BuildingID int64     `json:"building_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UnitStore struct {
	db *sql.DB
}

func (s *UnitStore) GetAll(ctx context.Context, buildingID int64) ([]Unit, error) {
	query := `
		SELECT id, name, building_id, created_at, updated_at
		FROM units
		WHERE building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var u Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.BuildingID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}

func (s *UnitStore) GetByID(ctx context.Context, id int64) (*Unit, error) {
	query := `
		SELECT id, name, building_id, created_at, updated_at
		FROM units
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var u Unit
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.BuildingID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (s *UnitStore) GetByIdTx(ctx context.Context, tx *sql.Tx, id int64) (*Unit, error) {
	query := `
		SELECT id, name, building_id, created_at, updated_at
		FROM units
		WHERE id = ?
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var u Unit
	err := tx.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.BuildingID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UnitStore) GetAllByPeopleID(ctx context.Context, peopleID int64) ([]Unit, error) {
	query := `
	SELECT DISTINCT u.id, u.name, u.building_id
		FROM units u
		INNER JOIN leases l ON u.id = l.unit_id
		WHERE l.people_id = ? AND l.status = '1'
		ORDER BY u.name
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, peopleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var u Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.BuildingID); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}
func (s *UnitStore) Create(ctx context.Context, unit *Unit) error {
	query := `
		INSERT INTO units (name, building_id)
		VALUES (?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, unit.Name, unit.BuildingID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	unit.ID = id
	return nil
}

func (s *UnitStore) Update(ctx context.Context, unit *Unit) error {
	query := `
		UPDATE units
		SET name = ?, building_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, unit.Name, unit.BuildingID, unit.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UnitStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM units
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// get available units by building id
func (s *UnitStore) GetAvailableUnitsByBuildingID(ctx context.Context, buildingID int64,includeUnitID *int64) ([]Unit, error) {
	query := `
		SELECT DISTINCT u.id, u.name, u.building_id
		FROM units u
		WHERE u.building_id = ?
			AND u.id NOT IN (
				SELECT DISTINCT l.unit_id
				FROM leases l
				WHERE l.building_id = ? AND l.status = '1'
			)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var u Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.BuildingID); err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	// get unit by id if includeUnitID is not 0
	if includeUnitID != nil {
		unit, err := s.GetByID(ctx, *includeUnitID)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}

	return units, nil
}
