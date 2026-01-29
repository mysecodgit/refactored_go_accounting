package service

import (
	"context"
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type BuildingStore interface {
	GetAll(ctx context.Context) ([]store.Building, error)
	GetAllByUserID(ctx context.Context, userID int64) ([]store.Building, error)
	GetByID(ctx context.Context, id int64) (*store.Building, error)
	Create(ctx context.Context, tx *sql.Tx, building *store.Building) error
	Update(ctx context.Context, building *store.Building) error
	Delete(ctx context.Context, id int64) error
}

type BuildingService struct {
	db *sql.DB
	buildingStore BuildingStore
	userBuildingStore UserBuildingStore
}

func NewBuildingService(db *sql.DB, buildingStore BuildingStore, userBuildingStore UserBuildingStore) *BuildingService {
	return &BuildingService{db: db, buildingStore: buildingStore, userBuildingStore: userBuildingStore}
}

func (s *BuildingService) GetAll(ctx context.Context) ([]store.Building, error) {
	return s.buildingStore.GetAll(ctx)
}

func (s *BuildingService) GetAllByUserID(ctx context.Context, userID int64) ([]store.Building, error) {
	return s.buildingStore.GetAllByUserID(ctx, userID)
}

func (s *BuildingService) GetByID(ctx context.Context, id int64) (*store.Building, error) {
	return s.buildingStore.GetByID(ctx, id)
}

func (s *BuildingService) Create(ctx context.Context, building *store.Building, userID int64) error {
	return withTx(s.db,ctx, func(tx *sql.Tx) error {
		if err := s.buildingStore.Create(ctx, tx, building); err != nil {
			return err
		}

		if err := s.userBuildingStore.AssignBuildingTX(ctx, tx,userID, building.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *BuildingService) Update(ctx context.Context, building *store.Building) error {
	return s.buildingStore.Update(ctx, building)
}

func (s *BuildingService) Delete(ctx context.Context, id int64) error {
	return s.buildingStore.Delete(ctx, id)
}
