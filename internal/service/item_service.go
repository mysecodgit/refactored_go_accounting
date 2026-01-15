package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type ItemStore interface {
	GetAll(ctx context.Context, buildingID int64) ([]store.Item, error)
	GetByID(ctx context.Context, id int64) (*store.Item, error)
	Create(ctx context.Context, i *store.Item) error
	Update(ctx context.Context, i *store.Item) error
	Delete(ctx context.Context, id int64) error
}

type ItemService struct {
	store ItemStore
}

func NewItemService(store ItemStore) *ItemService {
	return &ItemService{store: store}
}

func (s *ItemService) GetAll(ctx context.Context, buildingID int64) ([]store.Item, error) {
	return s.store.GetAll(ctx, buildingID)
}

func (s *ItemService) GetByID(ctx context.Context, id int64) (*store.Item, error) {
	return s.store.GetByID(ctx, id)
}

func (s *ItemService) Create(ctx context.Context, i *store.Item) error {
	return s.store.Create(ctx, i)
}

func (s *ItemService) Update(ctx context.Context, i *store.Item) error {
	return s.store.Update(ctx, i)
}

func (s *ItemService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
