package service

import (
	"context"
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type UserBuildingStore interface {
	GetBuildingsByUserID(ctx context.Context, userID int64) ([]store.Building, error)
	AssignBuilding(ctx context.Context, userID, buildingID int64) error
	AssignBuildingTX(ctx context.Context, tx *sql.Tx, userID, buildingID int64) error
	UnassignBuilding(ctx context.Context, userID, buildingID int64) error
	GetUsersByBuildingID(ctx context.Context, buildingID int64) ([]store.User, error)
}

type UserBuildingService struct {
	userBuildingStore UserBuildingStore
}

func NewUserBuildingService(userBuildingStore UserBuildingStore) *UserBuildingService {
	return &UserBuildingService{
		userBuildingStore: userBuildingStore,
	}
}

func (s *UserBuildingService) GetBuildingsByUserID(ctx context.Context, userID int64) ([]store.Building, error) {
	return s.userBuildingStore.GetBuildingsByUserID(ctx, userID)
}

func (s *UserBuildingService) AssignBuilding(ctx context.Context, userID, buildingID int64) error {
	return s.userBuildingStore.AssignBuilding(ctx, userID, buildingID)
}

func (s *UserBuildingService) UnassignBuilding(ctx context.Context, userID, buildingID int64) error {
	return s.userBuildingStore.UnassignBuilding(ctx, userID, buildingID)
}

func (s *UserBuildingService) GetUsersByBuildingID(ctx context.Context, buildingID int64) ([]store.User, error) {
	return s.userBuildingStore.GetUsersByBuildingID(ctx, buildingID)
}
