package service

import (
	"context"
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type UnitStore interface {
	GetAll(ctx context.Context, buildingID int64) ([]store.Unit, error)
	GetByID(ctx context.Context, id int64) (*store.Unit, error)
	GetByIdTx(ctx context.Context, tx *sql.Tx, id int64) (*store.Unit, error)
	GetAllByPeopleID(ctx context.Context, peopleID int64) ([]store.Unit, error)
	Create(ctx context.Context, unit *store.Unit) error
	Update(ctx context.Context, unit *store.Unit) error
	Delete(ctx context.Context, id int64) error
	GetAvailableUnitsByBuildingID(ctx context.Context, buildingID int64,includeUnitID *int64) ([]store.Unit, error)
}

type UnitService struct {
	unitStore UnitStore
}

func NewUnitService(unitStore UnitStore) *UnitService {
	return &UnitService{unitStore: unitStore}
}

func (s *UnitService) GetAll(ctx context.Context, buildingID int64) ([]store.Unit, error) {
	return s.unitStore.GetAll(ctx, buildingID)
}

func (s *UnitService) GetByID(ctx context.Context, id int64) (*store.Unit, error) {
	return s.unitStore.GetByID(ctx, id)
}

func (s *UnitService) GetAllByPeopleID(ctx context.Context, peopleID int64) ([]store.Unit, error) {
	return s.unitStore.GetAllByPeopleID(ctx, peopleID)
}

func (s *UnitService) Create(ctx context.Context, unit *store.Unit) error {
	return s.unitStore.Create(ctx, unit)
}

func (s *UnitService) Update(ctx context.Context, unit *store.Unit) error {
	return s.unitStore.Update(ctx, unit)
}

func (s *UnitService) Delete(ctx context.Context, id int64) error {
	return s.unitStore.Delete(ctx, id)
}
// get available units by building id
func (s *UnitService) GetAvailableUnitsByBuildingID(ctx context.Context, buildingID int64,includeUnitID *int64) ([]store.Unit, error) {
	return s.unitStore.GetAvailableUnitsByBuildingID(ctx, buildingID,includeUnitID)
}