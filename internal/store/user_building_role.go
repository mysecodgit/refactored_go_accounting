package store

import (
	"context"
	"database/sql"
)

type UserBuildingRole struct {
	UserID     int64 `json:"user_id"`
	BuildingID int64 `json:"building_id"`
	RoleID     int64 `json:"role_id"`
}

type UserBuildingRoleStore struct {
	db *sql.DB
}

func (s *UserBuildingRoleStore) GetRoleByUserAndBuilding(ctx context.Context, userID, buildingID int64) (*Role, error) {
	query := `
		SELECT r.id, r.owner_user_id, r.name, r.created_at
		FROM roles r
		INNER JOIN user_building_roles ubr ON r.id = ubr.role_id
		WHERE ubr.user_id = ? AND ubr.building_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var r Role
	err := s.db.QueryRowContext(ctx, query, userID, buildingID).Scan(&r.ID, &r.OwnerUserID, &r.Name, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &r, nil
}

func (s *UserBuildingRoleStore) GetUsersByBuildingAndRole(ctx context.Context, buildingID, roleID int64) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.username, u.phone, u.parent_user_id
		FROM users u
		INNER JOIN user_building_roles ubr ON u.id = ubr.user_id
		WHERE ubr.building_id = ? AND ubr.role_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, buildingID, roleID)
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

func (s *UserBuildingRoleStore) AssignRole(ctx context.Context, userID, buildingID, roleID int64) error {
	// Check if already assigned
	checkQuery := `
		SELECT COUNT(*) FROM user_building_roles
		WHERE user_id = ? AND building_id = ? AND role_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, checkQuery, userID, buildingID, roleID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Already assigned, return success
		return nil
	}

	query := `
		INSERT INTO user_building_roles (user_id, building_id, role_id)
		VALUES (?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query, userID, buildingID, roleID)
	return err
}

func (s *UserBuildingRoleStore) UnassignRole(ctx context.Context, userID, buildingID, roleID int64) error {
	query := `
		DELETE FROM user_building_roles
		WHERE user_id = ? AND building_id = ? AND role_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, userID, buildingID, roleID)
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

func (s *UserBuildingRoleStore) GetRolesByUserAndBuilding(ctx context.Context, userID, buildingID int64) ([]Role, error) {
	query := `
		SELECT r.id, r.owner_user_id, r.name, r.created_at
		FROM roles r
		INNER JOIN user_building_roles ubr ON r.id = ubr.role_id
		WHERE ubr.user_id = ? AND ubr.building_id = ?
		ORDER BY r.name
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.OwnerUserID, &r.Name, &r.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}
