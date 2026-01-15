package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type PeopleStore interface {
	GetAll(ctx context.Context, buildingID int64) ([]store.People, error)
	GetByID(ctx context.Context, id int64) (*store.People, error)
	Create(ctx context.Context, p *store.People) error
	Update(ctx context.Context, p *store.People) error
	Delete(ctx context.Context, id int64) error
}

type PeopleService struct {
	store PeopleStore
}

func NewPeopleService(store PeopleStore) *PeopleService {
	return &PeopleService{store: store}
}

func (s *PeopleService) GetAll(ctx context.Context, buildingID int64) ([]store.People, error) {
	return s.store.GetAll(ctx, buildingID)
}

func (s *PeopleService) GetByID(ctx context.Context, id int64) (*store.People, error) {
	return s.store.GetByID(ctx, id)
}

func (s *PeopleService) Create(ctx context.Context, p *store.People) error {
	return s.store.Create(ctx, p)
}

func (s *PeopleService) Update(ctx context.Context, p *store.People) error {
	return s.store.Update(ctx, p)
}

func (s *PeopleService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
