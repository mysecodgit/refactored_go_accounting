package store

import (
	"context"
	"database/sql"
)

type RolePermission struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

type RolePermissionStore struct {
	db *sql.DB
}

func (s *RolePermissionStore) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]Permission, error) {
	query := `
		SELECT p.id, p.module, p.action, p.` + "`key`" + `
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.module, p.action
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Module, &p.Action, &p.Key); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

func (s *RolePermissionStore) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	// Check if already assigned
	checkQuery := `
		SELECT COUNT(*) FROM role_permissions
		WHERE role_id = ? AND permission_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, checkQuery, roleID, permissionID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Already assigned
		return nil
	}

	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES (?, ?)
	`

	_, err = s.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

func (s *RolePermissionStore) UnassignPermission(ctx context.Context, roleID, permissionID int64) error {
	query := `
		DELETE FROM role_permissions
		WHERE role_id = ? AND permission_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, roleID, permissionID)
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

func (s *RolePermissionStore) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all existing permissions for this role
	deleteQuery := `DELETE FROM role_permissions WHERE role_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, roleID)
	if err != nil {
		return err
	}

	// Insert new permissions
	insertQuery := `INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)`
	for _, permID := range permissionIDs {
		_, err = tx.ExecContext(ctx, insertQuery, roleID, permID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
