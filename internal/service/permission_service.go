package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type PermissionStore interface {
	GetAll(ctx context.Context) ([]store.Permission, error)
	GetByID(ctx context.Context, id int64) (*store.Permission, error)
	GetByKey(ctx context.Context, key string) (*store.Permission, error)
	Create(ctx context.Context, permission *store.Permission) error
	Update(ctx context.Context, permission *store.Permission) error
	Delete(ctx context.Context, id int64) error
}

type PermissionService struct {
	permissionStore PermissionStore
}

func NewPermissionService(permissionStore PermissionStore) *PermissionService {
	return &PermissionService{
		permissionStore: permissionStore,
	}
}

func (s *PermissionService) GetAll(ctx context.Context) ([]store.Permission, error) {
	return s.permissionStore.GetAll(ctx)
}

func (s *PermissionService) GetByID(ctx context.Context, id int64) (*store.Permission, error) {
	return s.permissionStore.GetByID(ctx, id)
}

func (s *PermissionService) GetByKey(ctx context.Context, key string) (*store.Permission, error) {
	return s.permissionStore.GetByKey(ctx, key)
}

func (s *PermissionService) Create(ctx context.Context, permission *store.Permission) error {
	return s.permissionStore.Create(ctx, permission)
}

func (s *PermissionService) Update(ctx context.Context, permission *store.Permission) error {
	return s.permissionStore.Update(ctx, permission)
}

func (s *PermissionService) Delete(ctx context.Context, id int64) error {
	return s.permissionStore.Delete(ctx, id)
}
