package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type PeopleTypeStore interface {
	GetAll(ctx context.Context) ([]store.PeopleType, error)
	GetByID(ctx context.Context, id int64) (*store.PeopleType, error)
	Create(ctx context.Context, pt *store.PeopleType) error
	Update(ctx context.Context, pt *store.PeopleType) error
	Delete(ctx context.Context, id int64) error
}

type PeopleTypeService struct {
	store PeopleTypeStore
}

func NewPeopleTypeService(store PeopleTypeStore) *PeopleTypeService {
	return &PeopleTypeService{store: store}
}

func (s *PeopleTypeService) GetAll(ctx context.Context) ([]store.PeopleType, error) {
	return s.store.GetAll(ctx)
}

func (s *PeopleTypeService) GetByID(ctx context.Context, id int64) (*store.PeopleType, error) {
	return s.store.GetByID(ctx, id)
}

func (s *PeopleTypeService) Create(ctx context.Context, pt *store.PeopleType) error {
	return s.store.Create(ctx, pt)
}

func (s *PeopleTypeService) Update(ctx context.Context, pt *store.PeopleType) error {
	return s.store.Update(ctx, pt)
}

func (s *PeopleTypeService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
