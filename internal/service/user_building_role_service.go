package service

import (
	"context"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type UserBuildingRoleStore interface {
	GetRoleByUserAndBuilding(ctx context.Context, userID, buildingID int64) (*store.Role, error)
	GetUsersByBuildingAndRole(ctx context.Context, buildingID, roleID int64) ([]store.User, error)
	AssignRole(ctx context.Context, userID, buildingID, roleID int64) error
	UnassignRole(ctx context.Context, userID, buildingID, roleID int64) error
	GetRolesByUserAndBuilding(ctx context.Context, userID, buildingID int64) ([]store.Role, error)
}

type UserBuildingRoleService struct {
	userBuildingRoleStore UserBuildingRoleStore
}

func NewUserBuildingRoleService(userBuildingRoleStore UserBuildingRoleStore) *UserBuildingRoleService {
	return &UserBuildingRoleService{
		userBuildingRoleStore: userBuildingRoleStore,
	}
}

func (s *UserBuildingRoleService) GetRoleByUserAndBuilding(ctx context.Context, userID, buildingID int64) (*store.Role, error) {
	return s.userBuildingRoleStore.GetRoleByUserAndBuilding(ctx, userID, buildingID)
}

func (s *UserBuildingRoleService) GetUsersByBuildingAndRole(ctx context.Context, buildingID, roleID int64) ([]store.User, error) {
	return s.userBuildingRoleStore.GetUsersByBuildingAndRole(ctx, buildingID, roleID)
}

func (s *UserBuildingRoleService) AssignRole(ctx context.Context, userID, buildingID, roleID int64) error {
	return s.userBuildingRoleStore.AssignRole(ctx, userID, buildingID, roleID)
}

func (s *UserBuildingRoleService) UnassignRole(ctx context.Context, userID, buildingID, roleID int64) error {
	return s.userBuildingRoleStore.UnassignRole(ctx, userID, buildingID, roleID)
}

func (s *UserBuildingRoleService) GetRolesByUserAndBuilding(ctx context.Context, userID, buildingID int64) ([]store.Role, error) {
	return s.userBuildingRoleStore.GetRolesByUserAndBuilding(ctx, userID, buildingID)
}
