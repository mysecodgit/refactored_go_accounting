package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type ReadingStore interface {
	GetAll(ctx context.Context, buildingID int64, status *string) ([]store.Reading, error)
	GetByID(ctx context.Context, id int64) (*store.Reading, error)
	GetAllByUnitID(ctx context.Context, unitID int64) ([]store.ReadingByUnitResponse, error)
	GetLatest(ctx context.Context, itemID int64, unitID int64) (*store.Reading, error)
	GetLatestTx(ctx context.Context, tx *sql.Tx, itemID int64, unitID int64) (*store.Reading, error)
	Create(ctx context.Context, tx *sql.Tx, reading *store.Reading) error
	Update(ctx context.Context, reading *store.Reading) error
	Delete(ctx context.Context, id int64) error
}

type ReadingService struct {
	readingStore ReadingStore
	db           *sql.DB
}

func NewReadingService(readingStore ReadingStore, db *sql.DB) *ReadingService {
	return &ReadingService{readingStore: readingStore, db: db}
}

func (s *ReadingService) GetAll(ctx context.Context, buildingID int64, status *string) ([]store.Reading, error) {
	return s.readingStore.GetAll(ctx, buildingID, status)
}

func (s *ReadingService) GetByID(ctx context.Context, id int64) (*store.Reading, error) {
	return s.readingStore.GetByID(ctx, id)
}

func (s *ReadingService) GetAllByUnitID(ctx context.Context, unitID int64) ([]store.ReadingByUnitResponse, error) {
	return s.readingStore.GetAllByUnitID(ctx, unitID)
}

func (s *ReadingService) GetLatest(ctx context.Context, itemID int64, unitID int64) (*store.Reading, error) {
	return s.readingStore.GetLatest(ctx, itemID, unitID)
}

func (s *ReadingService) Create(ctx context.Context, req dto.CreateReadingRequest) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// add this later to the database
		// 		CREATE UNIQUE INDEX uniq_reading
		// ON readings (
		//     item_id,
		//     unit_id,
		//     reading_month,
		//     reading_year,
		//     previous_value,
		//     current_value
		// );
		for _, reading := range req.Readings {
			// check if the readings are already created
			latest, err := s.readingStore.GetLatestTx(ctx, tx, int64(reading.ItemID), int64(reading.UnitID))
			if err != nil {
				return err
			}

			if latest != nil {
				fmt.Println("reading.ReadingDate", reading.ReadingDate)
				fmt.Println("latest.ReadingDate", strings.Split(latest.ReadingDate, "T")[0])
				if reading.ReadingDate < strings.Split(latest.ReadingDate, "T")[0] {
					return errors.New("reading date must be after last reading date")
				}

				fmt.Println("reading.PreviousValue", *reading.PreviousValue)
				fmt.Println("latest.PreviousValue", *latest.PreviousValue)
				fmt.Println("reading.CurrentValue", *reading.CurrentValue)
				fmt.Println("latest.CurrentValue", *latest.CurrentValue)

				if *reading.PreviousValue == *latest.PreviousValue && *reading.CurrentValue == *latest.CurrentValue {
					return errors.New("this reading already exists")
				}

				fmt.Println("reading.PreviousValue", *reading.PreviousValue)
				fmt.Println("latest.CurrentValue", *latest.CurrentValue)
				if *reading.PreviousValue != *latest.CurrentValue {
					return errors.New("previous value does not match last current value")
				}
			}

			reading := &store.Reading{
				ItemID:        int64(reading.ItemID),
				UnitID:        int64(reading.UnitID),
				LeaseID:       reading.LeaseID,
				ReadingMonth:  reading.ReadingMonth,
				ReadingYear:   reading.ReadingYear,
				ReadingDate:   reading.ReadingDate,
				PreviousValue: reading.PreviousValue,
				CurrentValue:  reading.CurrentValue,
				UnitPrice:     reading.UnitPrice,
				TotalAmount:   reading.TotalAmount,
				Notes:         reading.Notes,
				Status:        reading.Status,
			}

			if err := s.readingStore.Create(ctx, tx, reading); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ReadingService) Update(ctx context.Context, req dto.UpdateReadingRequest) error {

	reading := &store.Reading{
		ID:            int64(req.ID),
		ItemID:        int64(req.ItemID),
		UnitID:        int64(req.UnitID),
		LeaseID:       req.LeaseID,
		ReadingMonth:  req.ReadingMonth,
		ReadingYear:   req.ReadingYear,
		ReadingDate:   req.ReadingDate,
		PreviousValue: req.PreviousValue,
		CurrentValue:  req.CurrentValue,
		UnitPrice:     req.UnitPrice,
		TotalAmount:   req.TotalAmount,
		Notes:         req.Notes,
		Status:        req.Status,
	}

	if err := s.readingStore.Update(ctx, reading); err != nil {
		return err
	}

	return nil

}

func (s *ReadingService) Delete(ctx context.Context, id int64) error {
	return s.readingStore.Delete(ctx, id)
}
