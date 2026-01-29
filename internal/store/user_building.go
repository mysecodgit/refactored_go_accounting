package store

import (
	"context"
	"database/sql"
)

type UserBuilding struct {
	UserID     int64 `json:"user_id"`
	BuildingID int64 `json:"building_id"`
}

type UserBuildingStore struct {
	db *sql.DB
}

func (s *UserBuildingStore) GetBuildingsByUserID(ctx context.Context, userID int64) ([]Building, error) {
	query := `
		SELECT b.id, b.name, b.created_at, b.updated_at
		FROM buildings b
		INNER JOIN users_building ub ON b.id = ub.building_id
		WHERE ub.user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []Building
	for rows.Next() {
		var b Building
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}

	return buildings, nil
}

func (s *UserBuildingStore) AssignBuilding(ctx context.Context, userID, buildingID int64) error {
	// Check if assignment already exists
	checkQuery := `
		SELECT COUNT(*) FROM users_building
		WHERE user_id = ? AND building_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, checkQuery, userID, buildingID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Already assigned, return success
		return nil
	}

	query := `
		INSERT INTO users_building (user_id, building_id)
		VALUES (?, ?)
	`

	_, err = s.db.ExecContext(ctx, query, userID, buildingID)
	return err
}

func (s *UserBuildingStore) UnassignBuilding(ctx context.Context, userID, buildingID int64) error {
	query := `
		DELETE FROM users_building
		WHERE user_id = ? AND building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, userID, buildingID)
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

func (s *UserBuildingStore) GetUsersByBuildingID(ctx context.Context, buildingID int64) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.username, u.phone, u.parent_user_id
		FROM users u
		INNER JOIN users_building ub ON u.id = ub.user_id
		WHERE ub.building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Phone, &u.ParentUserID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
