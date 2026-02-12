package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|---------------------------------------------------------------------------
| Interfaces
|---------------------------------------------------------------------------
*/

type BillStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]store.Bill, error)
	GetByID(ctx context.Context, id int64) (*store.Bill, error)
	Create(ctx context.Context, tx *sql.Tx, b *store.Bill) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, b *store.Bill) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

type BillExpenseLineStore interface {
	GetAllByBillID(ctx context.Context, billID int64) ([]store.BillExpenseLine, error)
	Create(ctx context.Context, tx *sql.Tx, l *store.BillExpenseLine) (*int64, error)
	DeleteByBillID(ctx context.Context, tx *sql.Tx, billID int64) error
}

/*
|---------------------------------------------------------------------------
| Service
|---------------------------------------------------------------------------
*/

type BillService struct {
	db                   *sql.DB
	billStore            BillStore
	billExpenseLineStore BillExpenseLineStore
	splitStore           SplitStore
	transactionStore     TransactionStore
	accountStore         AccountStore
}

/*
|---------------------------------------------------------------------------
| Constructor
|---------------------------------------------------------------------------
*/

func NewBillService(
	db *sql.DB,
	billStore BillStore,
	billExpenseLineStore BillExpenseLineStore,
	splitStore SplitStore,
	transactionStore TransactionStore,
	accountStore AccountStore,
) *BillService {
	return &BillService{
		db:                   db,
		billStore:            billStore,
		billExpenseLineStore: billExpenseLineStore,
		splitStore:           splitStore,
		transactionStore:     transactionStore,
		accountStore:         accountStore,
	}
}

/*
|---------------------------------------------------------------------------
| Queries
|---------------------------------------------------------------------------
*/

func (s *BillService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
	peopleID *int,
	status *string,
) ([]*dto.BillDto, error) {
	bills, err := s.billStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get bills: %v", err)
	}
	return dto.MapBillsToDtos(bills), nil
}

func (s *BillService) GetByID(ctx context.Context, id int64) (*dto.BillResponseDetails, error) {
	bill, err := s.billStore.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("bill not found: %v", err)
	}

	var dtoBill *dto.BillDto
	if bill != nil {
		dtoBill = dto.MapBillToDto(*bill)
	}

	expenseLines, err := s.billExpenseLineStore.GetAllByBillID(ctx, bill.ID)
	if err != nil {
		return nil, fmt.Errorf("expense lines not found: %v", err)
	}

	var dtoExpenseLines []*dto.BillExpenseLineDto
	if expenseLines != nil {
		dtoExpenseLines = dto.MapBillExpenseLinesToDtos(expenseLines)
	}

	splits, err := s.splitStore.GetByTransactionID(ctx, bill.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("splits not found: %v", err)
	}

	dtoSplits := dto.MapSplitsToDto(splits)

	transaction, err := s.transactionStore.GetByID(ctx, bill.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %v", err)
	}

	return &dto.BillResponseDetails{
		Bill:         *dtoBill,
		ExpenseLines: dtoExpenseLines,
		Splits:       dtoSplits,
		Transaction:  *transaction,
	}, nil
}

/*
|---------------------------------------------------------------------------
| Commands
|---------------------------------------------------------------------------
*/

func (s *BillService) Create(ctx context.Context, req dto.CreateBillRequest) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Create transaction
		transaction := &store.Transaction{
			Type:              "bill",
			TransactionDate:   req.BillDate,
			TransactionNumber: req.BillNo,
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            req.UnitID,
		}
		transactionID, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %v", err)
		}
		// Create bill
		bill := &store.Bill{
			TransactionID: *transactionID,
			BillNo:        req.BillNo,
			BillDate:      req.BillDate,
			DueDate:       req.DueDate,
			APAccountID:   req.APAccountID,
			UnitID:        req.UnitID,
			PeopleID:      req.PeopleID,
			UserID:        1, // TODO: get user id from jwt
			Amount:        req.Amount,
			AmountCents:   amountCents,
			Description:   req.Description,
			CancelReason:  nil,
			Status:        "1",
			BuildingID:    req.BuildingID,
		}
		billID, err := s.billStore.Create(ctx, tx, bill)
		if err != nil {
			return err
		}

		// Create expense lines
		for _, line := range req.ExpenseLines {

			amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
			amountCents, err := money.ParseUSDAmount(amountStr)
			if err != nil {
				return fmt.Errorf("failed to parse amount: %v", err)
			}

			expenseLine := &store.BillExpenseLine{
				BillID:      *billID,
				AccountID:   line.AccountID,
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Amount:      line.Amount,
				AmountCents: amountCents,
			}
			_, err = s.billExpenseLineStore.Create(ctx, tx, expenseLine)
			if err != nil {
				return err
			}
		}

		// Generate splits
		splits, err := s.GenerateBillSplits(ctx, req.BillPayloadDTO)
		if err != nil {
			return err
		}

		if err := validateSplitsBalanced(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionID
			err = s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BillService) Update(ctx context.Context, req dto.UpdateBillRequest, billID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Fetch existing bill
		existingBill, err := s.billStore.GetByID(ctx, billID)
		if err != nil {
			return fmt.Errorf("bill not found: %v", err)
		}

		// Delete expense lines
		expenseLines, err := s.billExpenseLineStore.GetAllByBillID(ctx, existingBill.ID)
		if err != nil {
			return fmt.Errorf("failed to get expense lines: %v", err)
		}

		if len(expenseLines) > 0 {
			if err := s.billExpenseLineStore.DeleteByBillID(ctx, tx, existingBill.ID); err != nil {
				return fmt.Errorf("failed to delete expense lines: %v", err)
			}
		}

		amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %v", err)
		}

		// Update bill
		updatedBill := &store.Bill{
			ID:            billID,
			TransactionID: existingBill.TransactionID,
			BillNo:        req.BillNo,
			BillDate:      req.BillDate,
			DueDate:       req.DueDate,
			APAccountID:   req.APAccountID,
			UnitID:        req.UnitID,
			PeopleID:      req.PeopleID,
			UserID:        1, // TODO: get user id from jwt
			Amount:        req.Amount,
			AmountCents:   amountCents,
			Description:   req.Description,
			CancelReason:  existingBill.CancelReason,
			Status:        existingBill.Status,
			BuildingID:    req.BuildingID,
		}

		_, err = s.billStore.Update(ctx, tx, updatedBill)
		if err != nil {
			return fmt.Errorf("failed to update bill: %v", err)
		}

		// Recreate expense lines
		for _, line := range req.ExpenseLines {
			amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
			amountCents, err := money.ParseUSDAmount(amountStr)
			if err != nil {
				return fmt.Errorf("failed to parse amount: %v", err)
			}

			expenseLine := &store.BillExpenseLine{
				BillID:      billID,
				AccountID:   line.AccountID,
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Amount:      line.Amount,
				AmountCents: amountCents,
			}
			_, err = s.billExpenseLineStore.Create(ctx, tx, expenseLine)
			if err != nil {
				return err
			}
		}

		// Update transaction
		transaction := &store.Transaction{
			ID:                existingBill.TransactionID,
			Type:              "bill",
			TransactionDate:   req.BillDate,
			TransactionNumber: req.BillNo,
			Memo:              req.Description,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            req.UnitID,
		}
		_, err = s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			return fmt.Errorf("failed to update transaction: %v", err)
		}

		// Soft delete splits
		if err := s.splitStore.DeleteByTransactionID(ctx, tx, existingBill.TransactionID); err != nil {
			return fmt.Errorf("failed to soft delete splits: %v", err)
		}

		// Generate splits
		splits, err := s.GenerateBillSplits(ctx, req.BillPayloadDTO)
		if err != nil {
			return err
		}

		if err := validateSplitsBalanced(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = existingBill.TransactionID
			err = s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *BillService) GenerateBillSplits(
	ctx context.Context,
	req dto.BillPayloadDTO,
) ([]store.Split, error) {
	// Validate AP account
	if _, err := s.accountStore.GetByID(ctx, req.APAccountID); err != nil {
		return nil, fmt.Errorf("AP account not found")
	}

	amount := req.Amount
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)
	amountCents, err := money.ParseUSDAmount(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %v", err)
	}

	splits := make([]store.Split, 0)

	// Credit AP account (liability increases)
	creditSplit := store.Split{
		AccountID:   req.APAccountID,
		Credit:      &amount,
		CreditCents: &amountCents,
		Debit:       nil,
		DebitCents:  nil,
		UnitID:      req.UnitID,
		PeopleID:    req.PeopleID,
		Status:      "1",
	}
	splits = append(splits, creditSplit)

	// Debit expense accounts
	for _, line := range req.ExpenseLines {
		_, err := s.accountStore.GetByID(ctx, line.AccountID)
		if err != nil {
			return nil, fmt.Errorf("account not found: %v", err)
		}

		amountStr := strconv.FormatFloat(line.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %v", err)
		}

		splits = append(splits, store.Split{
			AccountID:   line.AccountID,
			Debit:       &line.Amount,
			DebitCents:  &amountCents,
			Credit:      nil,
			CreditCents: nil,
			UnitID:      line.UnitID,
			PeopleID:    line.PeopleID,
			Status:      "1",
		})
	}

	return splits, nil
}

func validateSplitsBalanced(splits []store.Split) error {
	var debit, credit float64
	var debitCents, creditCents int64
	for _, s := range splits {
		if s.Debit != nil {
			debit += *s.Debit
		}
		if s.Credit != nil {
			credit += *s.Credit
		}
		if s.DebitCents != nil {
			debitCents += *s.DebitCents
		}
		if s.CreditCents != nil {
			creditCents += *s.CreditCents
		}
	}

	// if math.Abs(debit-credit) > 0.0001 {
	// 	return fmt.Errorf("unbalanced entry: debit %.2f ≠ credit %.2f", debit, credit)
	// }

	if debitCents != creditCents {
		return fmt.Errorf("unbalanced entry: debit cents %d ≠ credit cents %d", debitCents, creditCents)
	}

	return nil
}
