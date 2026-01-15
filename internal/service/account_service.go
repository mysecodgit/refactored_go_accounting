package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

// Interface for dependency injection
type AccountStore interface {
	GetAll(ctx context.Context, buildingID int64) ([]store.Account, error)
	GetByID(ctx context.Context, id int64) (*store.Account, error)
	Create(ctx context.Context, a *store.Account) error
	Update(ctx context.Context, a *store.Account) error
	Delete(ctx context.Context, id int64) error
}

type AccountService struct {
	store AccountStore
}

func NewAccountService(store AccountStore) *AccountService {
	return &AccountService{store: store}
}

func (s *AccountService) GetAll(ctx context.Context, buildingID int64) ([]store.Account, error) {
	return s.store.GetAll(ctx, buildingID)
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (*store.Account, error) {
	return s.store.GetByID(ctx, id)
}

func (s *AccountService) Create(ctx context.Context, a *store.Account) error {
	return s.store.Create(ctx, a)
}

func (s *AccountService) Update(ctx context.Context, a *store.Account) error {
	return s.store.Update(ctx, a)
}

func (s *AccountService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
