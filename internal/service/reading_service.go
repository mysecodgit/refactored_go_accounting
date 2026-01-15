package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type ReadingStore interface {
	GetAll(ctx context.Context,buildingID int64) ([]store.Reading, error)
	GetByID(ctx context.Context, id int64) (*store.Reading, error)
	GetAllByUnitID(ctx context.Context, unitID int64) ([]store.ReadingByUnitResponse, error)
	Create(ctx context.Context, reading *store.Reading) error
	Update(ctx context.Context, reading *store.Reading) error
	Delete(ctx context.Context, id int64) error
}

type ReadingService struct {
	readingStore ReadingStore
}

func NewReadingService(readingStore ReadingStore) *ReadingService {
	return &ReadingService{readingStore: readingStore}
}

func (s *ReadingService) GetAll(ctx context.Context,buildingID int64) ([]store.Reading, error) {
	return s.readingStore.GetAll(ctx,buildingID)
}

func (s *ReadingService) GetByID(ctx context.Context, id int64) (*store.Reading, error) {
	return s.readingStore.GetByID(ctx, id)
}

func (s *ReadingService) GetAllByUnitID(ctx context.Context, unitID int64) ([]store.ReadingByUnitResponse, error) {
	return s.readingStore.GetAllByUnitID(ctx, unitID)
}

func (s *ReadingService) Create(ctx context.Context, reading *store.Reading) error {
	return s.readingStore.Create(ctx, reading)
}

func (s *ReadingService) Update(ctx context.Context, reading *store.Reading) error {
	return s.readingStore.Update(ctx, reading)
}

func (s *ReadingService) Delete(ctx context.Context, id int64) error {
	return s.readingStore.Delete(ctx, id)
}
