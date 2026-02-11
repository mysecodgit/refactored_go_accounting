package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/dto"
	_ "github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|---------------------------------------------------------------------------
| Interfaces
|---------------------------------------------------------------------------
*/

type CheckStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string) ([]store.Check, error)
	GetByID(ctx context.Context, id int64) (*store.Check, error)
	Create(ctx context.Context, tx *sql.Tx, c *store.Check) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, c *store.Check) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

/*
|---------------------------------------------------------------------------
| Service
|---------------------------------------------------------------------------
*/

type CheckService struct {
	db               *sql.DB
	checkStore       CheckStore
	expenseLineStore ExpenseLineStore
	splitStore       SplitStore
	transactionStore TransactionStore
	accountStore     AccountStore
}

type ExpenseLineStore interface {
	GetAllByCheckID(ctx context.Context, checkID int64) ([]store.ExpenseLine, error)
	GetByID(ctx context.Context, id int64) (*store.ExpenseLine, error)
	Create(ctx context.Context, tx *sql.Tx, l *store.ExpenseLine) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, l *store.ExpenseLine) (*int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteByCheckID(ctx context.Context, tx *sql.Tx, checkID int64) error
}

/*
|---------------------------------------------------------------------------
| Constructor
|---------------------------------------------------------------------------
*/

func NewCheckService(db *sql.DB,
	checkStore CheckStore,
	expenseLineStore ExpenseLineStore,
	splitStore SplitStore,
	transactionStore TransactionStore,
	accountStore AccountStore,
) *CheckService {
	return &CheckService{
		db:               db,
		checkStore:       checkStore,
		expenseLineStore: expenseLineStore,
		splitStore:       splitStore,
		transactionStore: transactionStore,
		accountStore:     accountStore,
	}
}

/*
|---------------------------------------------------------------------------
| Queries
|---------------------------------------------------------------------------
*/

func (s *CheckService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate, endDate *string,
) ([]*dto.CheckDto, error) {
	checks, err := s.checkStore.GetAll(ctx, buildingID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get checks: %v", err)
	}
	return dto.MapChecksToDtos(checks), nil
}

func (s *CheckService) GetByID(ctx context.Context, id int64) (*dto.CheckResponseDetails, error) {

	check, err := s.checkStore.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("check not found: %v", err)
	}

	checkDto := dto.MapCheckToDto(*check)

	expenseLines, err := s.expenseLineStore.GetAllByCheckID(ctx, check.ID)
	if err != nil {
		return nil, fmt.Errorf("expense lines not found: %v", err)
	}

	expenseLinesDto := dto.MapExpenseLinesToDtos(expenseLines)

	splits, err := s.splitStore.GetByTransactionID(ctx, check.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("splits not found: %v", err)
	}

	splitsDto := dto.MapSplitsToDto(splits)

	transaction, err := s.transactionStore.GetByID(ctx, check.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %v", err)
	}

	return &dto.CheckResponseDetails{
		Check:        checkDto,
		ExpenseLines: expenseLinesDto,
		Splits:       splitsDto,
		Transaction:  *transaction,
	}, nil
}

/*
|---------------------------------------------------------------------------
| Commands
|---------------------------------------------------------------------------
*/
func (s *CheckService) Create(ctx context.Context, req dto.CreateCheckRequest) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create transaction
		transaction := &store.Transaction{
			Type:              "check",
			TransactionDate:   req.CheckDate,
			TransactionNumber: *req.ReferenceNumber,
			Memo:              *req.Memo,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            nil,
		}
		transactionId, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		totalAmountStr := strconv.FormatFloat(req.TotalAmount, 'f', -1, 64)
		totalAmountCents, err := money.ParseUSDAmount(totalAmountStr)
		if err != nil {
			return err
		}

		// create check
		check := &store.Check{
			TransactionID:    *transactionId,
			CheckDate:        req.CheckDate,
			ReferenceNumber:  *req.ReferenceNumber,
			PaymentAccountID: req.PaymentAccountID,
			BuildingID:       req.BuildingID,
			Memo:             req.Memo,
			TotalAmount:      req.TotalAmount,
			AmountCents:      totalAmountCents,
		}
		checkId, err := s.checkStore.Create(ctx, tx, check)
		if err != nil {
			return err
		}

		// create expense lines
		for _, line := range req.ExpenseLines {
			amountStr := strconv.FormatFloat(line.Amount, 'f', -1, 64)
			amountCents, err := money.ParseUSDAmount(amountStr)
			if err != nil {
				return err
			}
			expenseLine := &store.ExpenseLine{
				CheckID:     *checkId,
				AccountID:   line.AccountID,
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Amount:      line.Amount,
				AmountCents: amountCents,
			}
			_, err = s.expenseLineStore.Create(ctx, tx, expenseLine)
			if err != nil {
				return err
			}
		}

		// generate splits
		splits, err := s.GenerateCheckSplits(ctx, req.CheckPayloadDTO)
		if err != nil {
			return err
		}

		// validate splits
		if err := s.ValidateSplits(splits); err != nil {
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionId
			err = s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				return err
			}
		}
		return nil
	})

}

func (s *CheckService) Update(ctx context.Context, req dto.UpdateCheckRequest, checkId int64) error {
	fmt.Println("Update check", checkId)
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Fetch existing check
		existingCheck, err := s.checkStore.GetByID(ctx, checkId)
		if err != nil {
			fmt.Println("Check not found", err)
			return fmt.Errorf("check not found: %v", err)
		}

		fmt.Println("Existing check", existingCheck)

		expenseLines, err := s.expenseLineStore.GetAllByCheckID(ctx, existingCheck.ID)
		if err != nil {
			fmt.Println("Failed to get expense lines", err)
			return fmt.Errorf("failed to get expense lines: %v", err)
		}
		fmt.Println("Expense lines", expenseLines)

		// delete expense lines
		if len(expenseLines) > 0 {
			if err := s.expenseLineStore.DeleteByCheckID(ctx, tx, existingCheck.ID); err != nil {
				fmt.Println("Failed to delete expense lines", err)
				return fmt.Errorf("failed to delete expense lines: %v", err)
			}
		}

		amountStr := strconv.FormatFloat(req.TotalAmount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return err
		}

		// update check
		updatedCheck := &store.Check{
			ID:               checkId,
			TransactionID:    existingCheck.TransactionID,
			CheckDate:        req.CheckDate,
			ReferenceNumber:  *req.ReferenceNumber,
			PaymentAccountID: req.PaymentAccountID,
			BuildingID:       req.BuildingID,
			Memo:             req.Memo,
			TotalAmount:      req.TotalAmount,
			AmountCents:      amountCents,
		}

		_, err = s.checkStore.Update(ctx, tx, updatedCheck)
		if err != nil {
			fmt.Println("Failed to update check", err, updatedCheck)
			return fmt.Errorf("failed to update check: %v", err)
		}

		// recreate expense lines
		for _, line := range req.ExpenseLines {
			amountStr := strconv.FormatFloat(line.Amount, 'f', -1, 64)
			amountCents, err := money.ParseUSDAmount(amountStr)
			if err != nil {
				return err
			}

			expenseLine := &store.ExpenseLine{
				CheckID:     checkId,
				AccountID:   line.AccountID,
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Amount:      line.Amount,
				AmountCents: amountCents,
			}
			_, err = s.expenseLineStore.Create(ctx, tx, expenseLine)
			if err != nil {
				fmt.Println("Failed to create expense line", err)
				return err
			}
		}

		// update transaction
		transaction := &store.Transaction{
			ID:                existingCheck.TransactionID,
			Type:              "check",
			TransactionDate:   req.CheckDate,
			TransactionNumber: *req.ReferenceNumber,
			Memo:              *req.Memo,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            nil,
		}
		_, err = s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			fmt.Println("Failed to update transaction", err)
			return fmt.Errorf("failed to update transaction: %v", err)
		}

		// soft delete splits
		if err := s.splitStore.DeleteByTransactionID(ctx, tx, existingCheck.TransactionID); err != nil {
			fmt.Println("Failed to soft delete splits", err)
			return fmt.Errorf("failed to soft delete splits: %v", err)
		}

		// generate splits
		splits, err := s.GenerateCheckSplits(ctx, req.CheckPayloadDTO)
		if err != nil {
			fmt.Println("Failed to generate splits", err)
			return err
		}

		// validate splits
		if err := s.ValidateSplits(splits); err != nil {
			fmt.Println("Failed to validate splits", err)
			return err
		}

		for _, split := range splits {
			split.TransactionID = existingCheck.TransactionID
			err = s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				fmt.Println("Failed to create split", err)
				return err
			}
		}

		return nil
	})
}

/*





func (s *CheckService) Delete(ctx context.Context, id int64) error {
	return s.checkStore.Delete(ctx, id)
}
*/

func (s *CheckService) GenerateCheckSplits(
	ctx context.Context,
	req dto.CheckPayloadDTO,
) ([]store.Split, error) {

	// Validate accounts
	if _, err := s.accountStore.GetByID(ctx, int64(req.PaymentAccountID)); err != nil {
		return nil, fmt.Errorf("deposit account not found")
	}

	amount := req.TotalAmount

	amountStr := strconv.FormatFloat(amount, 'f', -1, 64)
	amountCents, err := money.ParseUSDAmount(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %v", err)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	splits := make([]store.Split, 0)

	debitSplit := store.Split{
		AccountID:   int64(req.PaymentAccountID),
		Debit:       &amount,
		DebitCents:  &amountCents,
		Credit:      nil,
		CreditCents: nil,
		UnitID:      nil,
		PeopleID:    nil,
		Status:      "1",
	}

	splits = append(splits, debitSplit)

	for _, line := range req.ExpenseLines {
		_, err := s.accountStore.GetByID(ctx, int64(line.AccountID))
		if err != nil {
			return nil, fmt.Errorf("account not found: %v", err)
		}

		amountStr := strconv.FormatFloat(line.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %v", err)
		}

		splits = append(splits, store.Split{
			AccountID:   int64(line.AccountID),
			Credit:      &amount,
			CreditCents: &amountCents,
			Debit:       nil,
			DebitCents:  nil,
			UnitID:      line.UnitID,
			PeopleID:    line.PeopleID,
			Status:      "1",
		})

	}

	return splits, nil
}

// validate splits
func (s *CheckService) ValidateSplits(splits []store.Split) error {
	totalDebitCents := int64(0)
	totalCreditCents := int64(0)
	for _, split := range splits {
		if split.DebitCents != nil {
			totalDebitCents += *split.DebitCents
		}
		if split.CreditCents != nil {
			totalCreditCents += *split.CreditCents
		}
	}

	debitAmount := money.FormatMoneyFromCents(totalDebitCents)
	creditAmount := money.FormatMoneyFromCents(totalCreditCents)

	if totalDebitCents != totalCreditCents {
		return fmt.Errorf("split must have %s debit and credit %s", debitAmount, creditAmount)
	}

	return nil
}
