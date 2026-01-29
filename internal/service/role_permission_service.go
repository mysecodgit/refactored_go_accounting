package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type RolePermissionStore interface {
	GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]store.Permission, error)
	AssignPermission(ctx context.Context, roleID, permissionID int64) error
	UnassignPermission(ctx context.Context, roleID, permissionID int64) error
	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}

type RolePermissionService struct {
	rolePermissionStore RolePermissionStore
}

func NewRolePermissionService(rolePermissionStore RolePermissionStore) *RolePermissionService {
	return &RolePermissionService{
		rolePermissionStore: rolePermissionStore,
	}
}

func (s *RolePermissionService) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]store.Permission, error) {
	return s.rolePermissionStore.GetPermissionsByRoleID(ctx, roleID)
}

func (s *RolePermissionService) AssignPermission(ctx context.Context, roleID, permissionID int64) error {
	return s.rolePermissionStore.AssignPermission(ctx, roleID, permissionID)
}

func (s *RolePermissionService) UnassignPermission(ctx context.Context, roleID, permissionID int64) error {
	return s.rolePermissionStore.UnassignPermission(ctx, roleID, permissionID)
}

func (s *RolePermissionService) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return s.rolePermissionStore.SetRolePermissions(ctx, roleID, permissionIDs)
}
