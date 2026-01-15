package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type AccountTypeStore interface {
	GetAll(ctx context.Context) ([]store.AccountType, error)
	GetByID(ctx context.Context, id int64) (*store.AccountType, error)
	Create(ctx context.Context, at *store.AccountType) error
	Update(ctx context.Context, at *store.AccountType) error
	Delete(ctx context.Context, id int64) error
}

type AccountTypeService struct {
	store AccountTypeStore
}

func NewAccountTypeService(store AccountTypeStore) *AccountTypeService {
	return &AccountTypeService{store: store}
}

func (s *AccountTypeService) GetAll(ctx context.Context) ([]store.AccountType, error) {
	return s.store.GetAll(ctx)
}

func (s *AccountTypeService) GetByID(ctx context.Context, id int64) (*store.AccountType, error) {
	return s.store.GetByID(ctx, id)
}

func (s *AccountTypeService) Create(ctx context.Context, at *store.AccountType) error {
	return s.store.Create(ctx, at)
}

func (s *AccountTypeService) Update(ctx context.Context, at *store.AccountType) error {
	return s.store.Update(ctx, at)
}

func (s *AccountTypeService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
