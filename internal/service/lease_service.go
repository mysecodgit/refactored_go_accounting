package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|--------------------------------------------------------------------------
| Store Interfaces
|--------------------------------------------------------------------------
*/

type LeaseStore interface {
	GetAll(
		ctx context.Context,
		buildingID int64,
		peopleID, unitID *int64,
		status *string,
	) ([]store.Lease, error)

	GetByID(ctx context.Context, id int64) (*store.Lease, error)

	GetActiveLeaseByUnitID(ctx context.Context, unitID int64) (*store.Lease, error)

	Create(ctx context.Context, tx *sql.Tx, lease *store.Lease) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, lease *store.Lease) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

type LeaseFileStore interface {
	GetByLeaseID(ctx context.Context, leaseID int64) ([]store.LeaseFile, error)
	Create(ctx context.Context, tx *sql.Tx, file *store.LeaseFile) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

/*
|--------------------------------------------------------------------------
| Service
|--------------------------------------------------------------------------
*/

type LeaseService struct {
	db             *sql.DB
	leaseStore     LeaseStore
	unitStore      UnitStore
	leaseFileStore LeaseFileStore
	peopleStore    PeopleStore
}

func NewLeaseService(
	db *sql.DB,
	leaseStore LeaseStore,
	unitStore UnitStore,
	leaseFileStore LeaseFileStore,
	peopleStore PeopleStore,
) *LeaseService {
	return &LeaseService{
		db:             db,
		leaseStore:     leaseStore,
		unitStore:      unitStore,
		leaseFileStore: leaseFileStore,
		peopleStore:    peopleStore,
	}
}

/*
|--------------------------------------------------------------------------
| Queries
|--------------------------------------------------------------------------
*/

func (s *LeaseService) GetAll(
	ctx context.Context,
	buildingID int64,
	peopleID, unitID *int64,
	status *string,
) ([]store.Lease, error) {
	return s.leaseStore.GetAll(ctx, buildingID, peopleID, unitID, status)
}

func (s *LeaseService) GetByID(ctx context.Context, id int64) (*dto.LeaseResponse, error) {

	lease, err := s.leaseStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	files, err := s.leaseFileStore.GetByLeaseID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.LeaseResponse{
		Lease:      *lease,
		LeaseFiles: files,
	}, nil
}

func (s *LeaseService) GetActiveLeaseByUnitID(ctx context.Context, unitID int64) ([]map[string]any, error) {
	lease, err := s.leaseStore.GetActiveLeaseByUnitID(ctx, unitID)
	if err != nil {
		return nil, err
	}

	people, err := s.peopleStore.GetByID(ctx, lease.PeopleID)
	if err != nil {
		return nil, err
	}
	return []map[string]any{
		{
			"lease":  lease,
			"people": people,
		},
	}, nil

}

/*
|--------------------------------------------------------------------------
| Create Lease + Upload Files
|--------------------------------------------------------------------------
*/

func (s *LeaseService) Create(
	ctx context.Context,
	req dto.CreateLeaseRequest,
	files []*multipart.FileHeader,
) (*dto.LeaseResponse, error) {

	var response dto.LeaseResponse
	var uploadedFiles []string

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {

		// get building id from unit id
		unit, err := s.unitStore.GetByID(ctx, int64(req.UnitID))
		if err != nil {
			return err
		}

		lease := &store.Lease{
			PeopleID:      int64(req.PeopleID),
			BuildingID:    unit.BuildingID,
			UnitID:        int64(req.UnitID),
			StartDate:     req.StartDate,
			EndDate:       req.EndDate,
			RentAmount:    req.RentAmount,
			DepositAmount: req.DepositAmount,
			ServiceAmount: req.ServiceAmount,
			LeaseTerms:    req.LeaseTerms,
			Status:        req.Status,
		}

		leaseID, err := s.leaseStore.Create(ctx, tx, lease)
		if err != nil {
			return err
		}

		lease.ID = *leaseID

		// upload files
		for _, file := range files {

			path, err := saveLeaseFile(*leaseID, file)
			if err != nil {
				return err
			}

			uploadedFiles = append(uploadedFiles, path)

			leaseFile := &store.LeaseFile{
				LeaseID:      *leaseID,
				Filename:     filepath.Base(path),
				OriginalName: file.Filename,
				FilePath:     path,
				FileType:     file.Header.Get("Content-Type"),
				FileSize:     file.Size,
			}

			if _, err := s.leaseFileStore.Create(ctx, tx, leaseFile); err != nil {
				return err
			}

			response.LeaseFiles = append(response.LeaseFiles, *leaseFile)
		}

		response.Lease = *lease
		return nil
	})

	// rollback files if tx fails
	if err != nil {
		for _, f := range uploadedFiles {
			_ = os.Remove(f)
		}
		return nil, err
	}

	return &response, nil
}

/*
|--------------------------------------------------------------------------
| Update Lease + Add Files
|--------------------------------------------------------------------------
*/

func (s *LeaseService) Update(
	ctx context.Context,
	req dto.UpdateLeaseRequest,
	files []*multipart.FileHeader,
) (*dto.LeaseResponse, error) {

	var response dto.LeaseResponse
	var uploadedFiles []string

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {

		existing, err := s.leaseStore.GetByID(ctx, int64(req.ID))
		if err != nil {
			fmt.Println("error getting lease by id", err)
			return err
		}

		// get building id from unit id
		unit, err := s.unitStore.GetByIdTx(ctx, tx, int64(req.UnitID))
		if err != nil {
			fmt.Println("error getting unit by id", err)
			return err
		}

		lease := &store.Lease{
			ID:            existing.ID,
			PeopleID:      int64(req.PeopleID),
			BuildingID:    unit.BuildingID,
			UnitID:        int64(req.UnitID),
			StartDate:     req.StartDate,
			EndDate:       req.EndDate,
			RentAmount:    req.RentAmount,
			DepositAmount: req.DepositAmount,
			ServiceAmount: req.ServiceAmount,
			LeaseTerms:    req.LeaseTerms,
			Status:        req.Status,
		}

		if _, err := s.leaseStore.Update(ctx, tx, lease); err != nil {
			fmt.Println("error updating lease", err)
			return err
		}

		// upload new files
		for _, file := range files {

			path, err := saveLeaseFile(existing.ID, file)
			if err != nil {
				fmt.Println("error saving lease file", err)
				return err
			}

			uploadedFiles = append(uploadedFiles, path)

			leaseFile := &store.LeaseFile{
				LeaseID:      existing.ID,
				Filename:     filepath.Base(path),
				OriginalName: file.Filename,
				FilePath:     path,
				FileType:     file.Header.Get("Content-Type"),
				FileSize:     file.Size,
			}

			if _, err := s.leaseFileStore.Create(ctx, tx, leaseFile); err != nil {
				fmt.Println("error creating lease file", err)
				return err
			}

			response.LeaseFiles = append(response.LeaseFiles, *leaseFile)
		}

		response.Lease = *existing
		return nil
	})

	if err != nil {
		fmt.Println("error updating lease", err)
		for _, f := range uploadedFiles {
			_ = os.Remove(f)
		}
		return nil, err
	}

	return &response, nil
}

/*
|--------------------------------------------------------------------------
| File Helper
|--------------------------------------------------------------------------
*/

func saveLeaseFile(leaseID int64, file *multipart.FileHeader) (string, error) {

	baseDir := filepath.Join("uploads", "leases", fmt.Sprintf("%d", leaseID))
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf(
		"%d_%d%s",
		leaseID,
		time.Now().UnixNano(),
		filepath.Ext(file.Filename),
	)

	dstPath := filepath.Join(baseDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return dstPath, nil
}
