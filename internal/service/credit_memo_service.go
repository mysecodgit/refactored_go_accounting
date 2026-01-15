package service

import (
	"context"
	"database/sql"
	"fmt"
	_ "fmt"
	"math"
	_ "math"

	"github.com/mysecodgit/go_accounting/internal/dto"
	_ "github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|--------------------------------------------------------------------------
| Interfaces
|--------------------------------------------------------------------------
*/

type CreditMemoStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]store.CreditMemoListResponse, error)
	GetByID(ctx context.Context, id int64) (*store.CreditMemo, error)
	Create(ctx context.Context, tx *sql.Tx, cm *store.CreditMemo) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, cm *store.CreditMemo) (*int64, error)
}

type CreditMemoService struct {
	db               *sql.DB
	creditMemoStore  CreditMemoStore
	transactionStore TransactionStore
	splitStore       SplitStore
	accountStore     AccountStore
}

/*
|--------------------------------------------------------------------------
| Constructor
|--------------------------------------------------------------------------
*/

func NewCreditMemoService(
	db *sql.DB,
	creditMemoStore CreditMemoStore,
	transactionStore TransactionStore,
	splitStore SplitStore,
	accountStore AccountStore,
) *CreditMemoService {
	return &CreditMemoService{
		db:               db,
		creditMemoStore:  creditMemoStore,
		transactionStore: transactionStore,
		splitStore:       splitStore,
		accountStore:     accountStore,
	}
}

/*
|--------------------------------------------------------------------------
| Queries
|--------------------------------------------------------------------------
*/

func (s *CreditMemoService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
	peopleID *int,
	status *string,
) ([]store.CreditMemoListResponse, error) {
	return s.creditMemoStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
}

type CreditMemoResponse struct {
	CreditMemo  *store.CreditMemo              `json:"credit_memo"`
	Splits      []store.Split          `json:"splits"`
	Transaction *store.Transaction `json:"transaction"`
}

func (s *CreditMemoService) GetByID(ctx context.Context, id int64) (*CreditMemoResponse, error) {
	creditMemo, err := s.creditMemoStore.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("credit memo not found: %v", err)
	}

	transaction, err := s.transactionStore.GetByID(ctx, creditMemo.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}

	allSplits, err := s.splitStore.GetByTransactionID(ctx, creditMemo.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch splits: %v", err)
	}

	return &CreditMemoResponse{
		CreditMemo:  creditMemo,
		Splits:      allSplits, // Include both active and inactive splits for view
		Transaction: transaction,
	}, nil
}



func (s *CreditMemoService) PreviewCreditMemoSplits(
	ctx context.Context,
	req dto.CreditMemoPayloadDTO,
) ([]store.Split, error) {

	splits, err := s.GenerateCreditMemoSplits(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := validateBalanced(splits); err != nil {
		return nil, err
	}

	return splits, nil
}


func (s *CreditMemoService) Create(
	ctx context.Context,
	req dto.CreateCreditMemoRequest,
) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// create transaction
		transaction := &store.Transaction{
			Type:              "credit memo",
			TransactionDate:   req.Date,
			TransactionNumber: req.Reference,
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: from JWT
			UnitID:            &req.UnitID,
		}

		transactionID, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// generate splits
		splits, err := s.GenerateCreditMemoSplits(ctx, req.CreditMemoPayloadDTO)
		if err != nil {
			return err
		}

		if err := s.validateBalanced(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				return err
			}
		}

		// create credit memo
		cm := &store.CreditMemo{
			TransactionID:    *transactionID,
			Reference:        req.Reference,
			Date:             req.Date,
			UserID:           1, // TODO: from JWT
			DepositTo:        req.DepositTo,
			LiabilityAccount: req.LiabilityAccount,
			PeopleID:         req.PeopleID,
			BuildingID:       req.BuildingID,
			UnitID:           req.UnitID,
			Amount:           req.Amount,
			Description:      req.Description,
			Status:           1,
		}

		_, err = s.creditMemoStore.Create(ctx, tx, cm)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *CreditMemoService) Update(
	ctx context.Context,
	req dto.UpdateCreditMemoRequest,
) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// 1️⃣ Get existing credit memo
		existingCM, err := s.creditMemoStore.GetByID(ctx, int64(req.ID))
		if err != nil {
			return fmt.Errorf("credit memo not found: %v", err)
		}

		// 2️⃣ Update transaction
		transaction := &store.Transaction{
			ID:                existingCM.TransactionID,
			Type:              "credit_memo",
			TransactionDate:   req.Date,
			TransactionNumber: req.Reference,
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get from JWT
			UnitID:            &req.UnitID,
		}

		transactionID, err := s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			return fmt.Errorf("error updating transaction: %v", err)
		}

		// 3️⃣ Delete existing splits
		if err := s.splitStore.DeleteByTransactionID(ctx, tx, existingCM.TransactionID); err != nil {
			return fmt.Errorf("error deleting existing splits: %v", err)
		}

		// 4️⃣ Generate new splits
		splits, err := s.GenerateCreditMemoSplits(ctx, req.CreditMemoPayloadDTO)
		if err != nil {
			return fmt.Errorf("error generating splits: %v", err)
		}

		// 5️⃣ Validate splits are balanced
		if err := validateBalanced(splits); err != nil {
			return fmt.Errorf("splits not balanced: %v", err)
		}

		// 6️⃣ Save new splits
		for _, split := range splits {
			split.TransactionID = *transactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				return fmt.Errorf("error creating split: %v", err)
			}
		}

		// 7️⃣ Update credit memo record
		updatedCM := &store.CreditMemo{
			ID:               existingCM.ID,
			TransactionID:    *transactionID,
			Reference:        req.Reference,
			Date:             req.Date,
			UserID:           1, // TODO: get from JWT
			DepositTo:        req.DepositTo,
			LiabilityAccount: req.LiabilityAccount,
			PeopleID:         req.PeopleID,
			BuildingID:       req.BuildingID,
			UnitID:           req.UnitID,
			Amount:           req.Amount,
			Description:      req.Description,
		}

		_, err = s.creditMemoStore.Update(ctx, tx, updatedCM)
		if err != nil {
			return fmt.Errorf("error updating credit memo: %v", err)
		}

		return nil
	})
}



func (s *CreditMemoService) GenerateCreditMemoSplits(
	ctx context.Context,
	req dto.CreditMemoPayloadDTO,
) ([]store.Split, error) {

	// Validate accounts
	if _, err := s.accountStore.GetByID(ctx, int64(req.DepositTo)); err != nil {
		return nil, fmt.Errorf("deposit account not found")
	}

	if _, err := s.accountStore.GetByID(ctx, int64(req.LiabilityAccount)); err != nil {
		return nil, fmt.Errorf("liability account not found")
	}

	amount := req.Amount
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// Credit Memo logic:
	// Debit: Liability account (reduces liability)
	// Credit: Deposit/Asset account

	debitSplit := store.Split{
		AccountID: int64(req.LiabilityAccount),
		Debit:     &amount,
		Credit:    nil,
		UnitID:    &req.UnitID,
		PeopleID:  &req.PeopleID,
		Status:    "1",
	}

	creditSplit := store.Split{
		AccountID: int64(req.DepositTo),
		Credit:    &amount,
		Debit:     nil,
		UnitID:    &req.UnitID,
		PeopleID:  &req.PeopleID,
		Status:    "1",
	}

	return []store.Split{debitSplit, creditSplit}, nil
}



func (s *CreditMemoService) validateBalanced(splits []store.Split) error {
	var debit, credit float64

	for _, s := range splits {
		if s.Debit != nil {
			debit += *s.Debit
		}
		if s.Credit != nil {
			credit += *s.Credit
		}
	}

	if math.Abs(debit-credit) > 0.0001 {
		return fmt.Errorf("unbalanced entry: debit %.2f ≠ credit %.2f", debit, credit)
	}

	return nil
}
