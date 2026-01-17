package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|--------------------------------------------------------------------------
| Store Interfaces
|--------------------------------------------------------------------------
*/

type SalesReceiptStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]store.SalesReceiptListResponse, error)
	GetByID(ctx context.Context, id int64) (*store.SalesReceipt, error)
	Create(ctx context.Context, tx *sql.Tx, receipt *store.SalesReceipt) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, receipt *store.SalesReceipt) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

type ReceiptItemStore interface {
	GetByReceiptID(ctx context.Context, receiptID int64) ([]store.ReceiptItem, error)
	Create(ctx context.Context, tx *sql.Tx, item *store.ReceiptItem) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, item *store.ReceiptItem) (*int64, error)
	DeleteByReceiptID(ctx context.Context, tx *sql.Tx, receiptID int64) error
}

/*
|--------------------------------------------------------------------------
| Service
|--------------------------------------------------------------------------
*/

type SalesReceiptService struct {
	db               *sql.DB
	salesReceiptStore SalesReceiptStore
	receiptItemStore  ReceiptItemStore
	transactionStore  TransactionStore
	splitStore        SplitStore
	accountStore      AccountStore
	itemStore         ItemStore
}

func NewSalesReceiptService(
	db *sql.DB,
	salesReceiptStore SalesReceiptStore,
	receiptItemStore ReceiptItemStore,
	transactionStore TransactionStore,
	splitStore SplitStore,
	accountStore AccountStore,
	itemStore ItemStore,
) *SalesReceiptService {
	return &SalesReceiptService{
		db:               db,
		salesReceiptStore: salesReceiptStore,
		receiptItemStore:  receiptItemStore,
		transactionStore:  transactionStore,
		splitStore:        splitStore,
		accountStore:      accountStore,
		itemStore:         itemStore,
	}
}

/*
|--------------------------------------------------------------------------
| Queries
|--------------------------------------------------------------------------
*/

func (s *SalesReceiptService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
	peopleID *int,
	status *string,
) ([]store.SalesReceiptListResponse, error) {
	return s.salesReceiptStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
}

func (s *SalesReceiptService) GetByID(ctx context.Context, id int64) (map[string]any, error) {
	receipt, err := s.salesReceiptStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.receiptItemStore.GetByReceiptID(ctx, id)
	if err != nil {
		return nil, err
	}

	splits, err := s.splitStore.GetAll(ctx, receipt.TransactionID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"receipt": receipt,
		"items":   items,
		"splits":  splits,
	}, nil
}

/*
|--------------------------------------------------------------------------
| Create
|--------------------------------------------------------------------------
*/

func (s *SalesReceiptService) Create(
	ctx context.Context,
	req dto.CreateSalesReceiptRequest,
) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// 1. Create transaction
		transaction := &store.Transaction{
			Type:              "receipt",
			TransactionDate:   req.ReceiptDate,
			TransactionNumber: fmt.Sprintf("%d", req.ReceiptNo),
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: JWT
			UnitID:            req.UnitID,
		}

		transactionID, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// 2. Generate splits
		splits, err := s.GenerateSalesReceiptSplits(ctx, req.SalesReceiptPayload)
		if err != nil {
			return err
		}

		if err := s.ValidateBalanced(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				return err
			}
		}

		// 3. Create receipt
		receipt := &store.SalesReceipt{
			TransactionID: *transactionID,
			ReceiptNo:     req.ReceiptNo,
			ReceiptDate:   req.ReceiptDate,
			UnitID:        req.UnitID,
			PeopleID:      req.PeopleID,
			UserID:        1,
			AccountID:     req.AccountID,
			Amount:        req.Amount,
			Description:   &req.Description,
			Status:        1,
			BuildingID:    req.BuildingID,
		}

		receiptID, err := s.salesReceiptStore.Create(ctx, tx, receipt)
		if err != nil {
			return err
		}

		// 4. Create receipt items
		for _, line := range req.Items {

			itemRow, err := s.itemStore.GetByID(ctx, int64(line.ItemID))
			if err != nil {
				return err
			}

			item := &store.ReceiptItem{
				ReceiptID:     *receiptID,
				ItemID:        int64(line.ItemID),
				ItemName:      itemRow.Name,
				Qty:           line.Qty,
				Rate:          line.Rate,
				Total:         *line.Total,
				PreviousValue: line.PreviousValue,
				CurrentValue:  line.CurrentValue,
			}

			if _, err := s.receiptItemStore.Create(ctx, tx, item); err != nil {
				return err
			}
		}

		return nil
	})
}

/*
|--------------------------------------------------------------------------
| Update
|--------------------------------------------------------------------------
*/

func (s *SalesReceiptService) Update(
	ctx context.Context,
	req dto.UpdateSalesReceiptRequest,
) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		existing, err := s.salesReceiptStore.GetByID(ctx, int64(req.ID))
		if err != nil {
			return err
		}

		// Update transaction
		transaction := &store.Transaction{
			ID:                existing.TransactionID,
			Type:              "receipt",
			TransactionDate:   req.ReceiptDate,
			TransactionNumber: fmt.Sprintf("%d", req.ReceiptNo),
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1,
			UnitID:            req.UnitID,
		}

		transactionID, err := s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// Rebuild splits
		if err := s.splitStore.DeleteByTransactionID(ctx, tx, existing.TransactionID); err != nil {
			return err
		}

		splits, err := s.GenerateSalesReceiptSplits(ctx, req.SalesReceiptPayload)
		if err != nil {
			return err
		}

		if err := s.ValidateBalanced(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				return err
			}
		}

		// Update receipt
		receipt := &store.SalesReceipt{
			ID:            existing.ID,
			TransactionID: *transactionID,
			ReceiptNo:     req.ReceiptNo,
			ReceiptDate:   req.ReceiptDate,
			UnitID:        req.UnitID,
			PeopleID:      req.PeopleID,
			AccountID:     req.AccountID,
			Amount:        req.Amount,
			Description:   &req.Description,
			BuildingID:    req.BuildingID,
			UserID:        1,
		}

		if _, err := s.salesReceiptStore.Update(ctx, tx, receipt); err != nil {
			return err
		}

		// Rebuild items
		if err := s.receiptItemStore.DeleteByReceiptID(ctx, tx, existing.ID); err != nil {
			return err
		}

		for _, line := range req.Items {

			itemRow, err := s.itemStore.GetByID(ctx, int64(line.ItemID))
			if err != nil {
				return err
			}

			item := &store.ReceiptItem{
				ReceiptID:     existing.ID,
				ItemID:        int64(line.ItemID),
				ItemName:      itemRow.Name,
				Qty:           line.Qty,
				Rate:          line.Rate,
				Total:         *line.Total,
				PreviousValue: line.PreviousValue,
				CurrentValue:  line.CurrentValue,
			}

			if _, err := s.receiptItemStore.Create(ctx, tx, item); err != nil {
				return err
			}
		}

		return nil
	})
}

/*
|--------------------------------------------------------------------------
| Split Generation
|--------------------------------------------------------------------------
*/

func (s *SalesReceiptService) GenerateSalesReceiptSplits(
	ctx context.Context,
	req dto.SalesReceiptPayload,
) ([]store.Split, error) {

	acc := make(map[int64]*splitAccumulator)

	addDebit := func(accountID int64, amount float64) {
		if acc[accountID] == nil {
			acc[accountID] = &splitAccumulator{}
		}
		acc[accountID].Debit += amount
	}

	addCredit := func(accountID int64, amount float64) {
		if acc[accountID] == nil {
			acc[accountID] = &splitAccumulator{}
		}
		acc[accountID].Credit += amount
	}

	// 1. Asset account debit
	addDebit(int64(req.AccountID), req.Amount)

	// 2. Item income
	for _, line := range req.Items {

		item, err := s.itemStore.GetByID(ctx, int64(line.ItemID))
		if err != nil {
			return nil, err
		}

		addCredit(*item.IncomeAccount, *line.Total)
	}

	// 3. Build splits
	var splits []store.Split

	for accountID, v := range acc {

		var debit, credit *float64

		if v.Debit > 0 {
			d := v.Debit
			debit = &d
		}
		if v.Credit > 0 {
			c := v.Credit
			credit = &c
		}

		splits = append(splits, store.Split{
			AccountID: accountID,
			Debit:     debit,
			Credit:    credit,
			UnitID:    req.UnitID,
			PeopleID:  req.PeopleID,
			Status:    "1",
		})
	}

	return splits, nil
}

/*
|--------------------------------------------------------------------------
| Validation
|--------------------------------------------------------------------------
*/

func (s *SalesReceiptService) ValidateBalanced(splits []store.Split) error {
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
