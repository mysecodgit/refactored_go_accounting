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

type BillPaymentStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate *string, endDate *string, peopleID *int, status *string) ([]store.BillPayment, error)
	GetAllByBillID(ctx context.Context, billID int64) ([]store.BillPayment, error)
	GetByID(ctx context.Context, id int64) (*store.BillPayment, error)
	GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*store.BillPayment, error)
	Create(ctx context.Context, tx *sql.Tx, billPayment *store.BillPayment) (*store.BillPayment, error)
	Update(ctx context.Context, tx *sql.Tx, billPayment *store.BillPayment) (*store.BillPayment, error)
	Delete(ctx context.Context, id int64) error
}

/*
|---------------------------------------------------------------------------
| Service
|---------------------------------------------------------------------------
*/

type BillPaymentService struct {
	db               *sql.DB
	billPaymentStore BillPaymentStore
	transactionStore TransactionStore
	accountStore     AccountStore
	billStore        BillStore
	splitStore       SplitStore
}

/*
|---------------------------------------------------------------------------
| Constructor
|---------------------------------------------------------------------------
*/

func NewBillPaymentService(
	db *sql.DB,
	billPaymentStore BillPaymentStore,
	transactionStore TransactionStore,
	accountStore AccountStore,
	billStore BillStore,
	splitStore SplitStore,
) *BillPaymentService {
	return &BillPaymentService{
		db:               db,
		billPaymentStore: billPaymentStore,
		transactionStore: transactionStore,
		accountStore:     accountStore,
		billStore:        billStore,
		splitStore:       splitStore,
	}
}

/*
|---------------------------------------------------------------------------
| Queries
|---------------------------------------------------------------------------
*/

func (s *BillPaymentService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate *string,
	endDate *string,
	peopleID *int,
	status *string,
) ([]*dto.BillPaymentDto, error) {
	payments, err := s.billPaymentStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		return nil, err
	}
	return dto.MapBillPaymentsToDtos(payments), nil
}

func (s *BillPaymentService) GetByID(
	ctx context.Context,
	id int64,
) (*dto.BillPaymentResponse, error) {
	payment, err := s.billPaymentStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dtoPayment := dto.MapBillPaymentToDto(*payment)

	splits, err := s.splitStore.GetByTransactionID(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	dtoSplits := dto.MapSplitsToDto(splits)

	transaction, err := s.transactionStore.GetByID(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	bill, err := s.billStore.GetByID(ctx, payment.BillID)
	if err != nil {
		return nil, err
	}

	dtoBill := dto.MapBillToDto(*bill)

	apAccount, err := s.accountStore.GetByID(ctx, bill.APAccountID)
	if err != nil {
		return nil, err
	}

	return &dto.BillPaymentResponse{
		Payment:     *dtoPayment,
		Splits:      dtoSplits,
		Transaction: *transaction,
		Bill:        *dtoBill,
		APAccount:   apAccount,
	}, nil
}

/*
|---------------------------------------------------------------------------
| Commands
|---------------------------------------------------------------------------
*/

func (s *BillPaymentService) Create(ctx context.Context, paymentDTO dto.CreateBillPaymentRequest) (*dto.BillPaymentResponse, error) {
	var response dto.BillPaymentResponse

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {
		bill, err := s.billStore.GetByID(ctx, int64(paymentDTO.BillID))
		if err != nil {
			return fmt.Errorf("bill not found: %v", err)
		}

		// Validate bill belongs to the building
		if bill.BuildingID != paymentDTO.BuildingID {
			return fmt.Errorf("bill does not belong to the specified building")
		}

		apAccount, err := s.accountStore.GetByID(ctx, bill.APAccountID)
		if err != nil {
			return fmt.Errorf("A/P account not found: %v", err)
		}

		// Get Asset Account from request
		assetAccount, err := s.accountStore.GetByID(ctx, int64(paymentDTO.AccountID))
		if err != nil {
			return fmt.Errorf("asset account not found: %v", err)
		}

		// Validate amount
		if paymentDTO.Amount == 0 {
			return fmt.Errorf("amount cannot be zero")
		}

		// Create transaction
		transaction := &store.Transaction{
			Type:              "bill_payment",
			TransactionDate:   paymentDTO.Date,
			TransactionNumber: paymentDTO.Reference,
			Memo:              paymentDTO.Reference,
			Status:            "1",
			BuildingID:        paymentDTO.BuildingID,
			UnitID:            bill.UnitID,
			UserID:            1, // TODO: get user id from jwt
		}
		transactionID, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		amountStr := strconv.FormatFloat(paymentDTO.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %v", err)
		}

		// Create splits
		// Debit Asset Account (cash/bank)
		debitSplit := store.Split{
			TransactionID: *transactionID,
			AccountID:     assetAccount.ID,
			Debit:         &paymentDTO.Amount,
			DebitCents:    &amountCents,
			Credit:        nil,
			CreditCents:   nil,
			UnitID:        bill.UnitID,
			PeopleID:      bill.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &debitSplit)
		if err != nil {
			return err
		}

		// Credit A/P Account (liability decreases)
		creditSplit := store.Split{
			TransactionID: *transactionID,
			AccountID:     apAccount.ID,
			Credit:        &paymentDTO.Amount,
			CreditCents:   &amountCents,
			Debit:         nil,
			DebitCents:   nil,
			UnitID:        bill.UnitID,
			PeopleID:      bill.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &creditSplit)
		if err != nil {
			return err
		}

		// Create bill payment
		billPayment := &store.BillPayment{
			TransactionID: *transactionID,
			Reference:     paymentDTO.Reference,
			Date:          paymentDTO.Date,
			BillID:        int64(paymentDTO.BillID),
			UserID:        1, // TODO: get user id from jwt
			AccountID:     int64(paymentDTO.AccountID),
			Amount:        paymentDTO.Amount,
			AmountCents:   amountCents,
			Status:        "1",
		}

		_, err = s.billPaymentStore.Create(ctx, tx, billPayment)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *BillPaymentService) Update(
	ctx context.Context,
	req dto.UpdateBillPaymentRequest,
	paymentID int64,
) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Fetch existing payment
		existing, err := s.billPaymentStore.GetByIDTx(ctx, tx, paymentID)
		if err != nil {
			return fmt.Errorf("bill payment not found: %v", err)
		}

		// Validate account
		if _, err := s.accountStore.GetByID(ctx, int64(req.AccountID)); err != nil {
			return fmt.Errorf("account not found")
		}

		amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %v", err)
		}

		// Update bill payment
		updatedPayment := &store.BillPayment{
			ID:            paymentID,
			TransactionID: existing.TransactionID,
			Reference:     req.Reference,
			Date:          req.Date,
			BillID:        existing.BillID,
			UserID:        1, // TODO: get user id from jwt
			AccountID:     int64(req.AccountID),
			Amount:        req.Amount,
			AmountCents:   amountCents,
			Status:        strconv.Itoa(req.Status),
		}

		if _, err := s.billPaymentStore.Update(ctx, tx, updatedPayment); err != nil {
			return err
		}

		// Update transaction
		transaction := &store.Transaction{
			ID:                existing.TransactionID,
			Type:              "bill_payment",
			TransactionDate:   req.Date,
			TransactionNumber: req.Reference,
			Memo:              req.Reference,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
		}
		_, err = s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// Soft delete splits
		if err := s.splitStore.DeleteByTransactionID(ctx, tx, existing.TransactionID); err != nil {
			return err
		}

		// Get bill for unit/people info
		bill, err := s.billStore.GetByID(ctx, existing.BillID)
		if err != nil {
			return err
		}

		// Get AP account
		apAccount, err := s.accountStore.GetByID(ctx, bill.APAccountID)
		if err != nil {
			return err
		}

		// Recreate splits
		// Debit Asset Account
		debitSplit := store.Split{
			TransactionID: existing.TransactionID,
			AccountID:     int64(req.AccountID),
			Debit:         &req.Amount,
			DebitCents:    &amountCents,
			Credit:        nil,
			CreditCents:   nil,
			UnitID:        bill.UnitID,
			PeopleID:      bill.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &debitSplit)
		if err != nil {
			return err
		}

		// Credit A/P Account
		creditSplit := store.Split{
			TransactionID: existing.TransactionID,
			AccountID:     apAccount.ID,
			Credit:        &req.Amount,
			CreditCents:   &amountCents,
			Debit:         nil,
			DebitCents:   nil,
			UnitID:        bill.UnitID,
			PeopleID:      bill.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &creditSplit)
		if err != nil {
			return err
		}

		return nil
	})
}
