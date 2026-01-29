package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type RoleStore interface {
	GetAllByOwnerID(ctx context.Context, ownerUserID int64) ([]store.Role, error)
	GetByID(ctx context.Context, id int64) (*store.Role, error)
	Create(ctx context.Context, role *store.Role) error
	Update(ctx context.Context, role *store.Role) error
	Delete(ctx context.Context, id int64, ownerUserID int64) error
}

type RoleService struct {
	roleStore RoleStore
}

func NewRoleService(roleStore RoleStore) *RoleService {
	return &RoleService{
		roleStore: roleStore,
	}
}

func (s *RoleService) GetAllByOwnerID(ctx context.Context, ownerUserID int64) ([]store.Role, error) {
	return s.roleStore.GetAllByOwnerID(ctx, ownerUserID)
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (*store.Role, error) {
	return s.roleStore.GetByID(ctx, id)
}

func (s *RoleService) Create(ctx context.Context, role *store.Role) error {
	return s.roleStore.Create(ctx, role)
}

func (s *RoleService) Update(ctx context.Context, role *store.Role) error {
	return s.roleStore.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id int64, ownerUserID int64) error {
	return s.roleStore.Delete(ctx, id, ownerUserID)
}
