package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type BuildingStore interface {
	GetAll(ctx context.Context) ([]store.Building, error)
	GetByID(ctx context.Context, id int64) (*store.Building, error)
	Create(ctx context.Context, building *store.Building) error
	Update(ctx context.Context, building *store.Building) error
	Delete(ctx context.Context, id int64) error
}

type BuildingService struct {
	buildingStore BuildingStore
}

func NewBuildingService(buildingStore BuildingStore) *BuildingService {
	return &BuildingService{buildingStore: buildingStore}
}

func (s *BuildingService) GetAll(ctx context.Context) ([]store.Building, error) {
	return s.buildingStore.GetAll(ctx)
}

func (s *BuildingService) GetByID(ctx context.Context, id int64) (*store.Building, error) {
	return s.buildingStore.GetByID(ctx, id)
}

func (s *BuildingService) Create(ctx context.Context, building *store.Building) error {
	return s.buildingStore.Create(ctx, building)
}

func (s *BuildingService) Update(ctx context.Context, building *store.Building) error {
	return s.buildingStore.Update(ctx, building)
}

func (s *BuildingService) Delete(ctx context.Context, id int64) error {
	return s.buildingStore.Delete(ctx, id)
}
