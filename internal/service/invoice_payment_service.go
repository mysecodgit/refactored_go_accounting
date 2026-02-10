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

type InvoicePaymentStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate *string, endDate *string, peopleID *int, status *string) ([]store.InvoicePayment, error)
	GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]store.InvoicePayment, error)
	GetByID(ctx context.Context, id int64) (*store.InvoicePayment, error)
	GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*store.InvoicePayment, error)
	Create(ctx context.Context, tx *sql.Tx, invoicePayment *store.InvoicePayment) (*store.InvoicePayment, error)
	Update(ctx context.Context, tx *sql.Tx, invoicePayment *store.InvoicePayment) (*store.InvoicePayment, error)
	Delete(ctx context.Context, id int64) error
}


/*
|---------------------------------------------------------------------------
| Service
|---------------------------------------------------------------------------
*/

type InvoicePaymentService struct {
	db                   *sql.DB
	invoicePaymentStore  InvoicePaymentStore
	transactionStore     TransactionStore
	accountStore         AccountStore
	invoiceStore         InvoiceStore
	splitStore           SplitStore
}

/*
|---------------------------------------------------------------------------
| Constructor
|---------------------------------------------------------------------------
*/

func NewInvoicePaymentService(
	db *sql.DB,
	invoicePaymentStore InvoicePaymentStore,
	transactionStore TransactionStore,
	accountStore AccountStore,
	invoiceStore InvoiceStore,
	splitStore SplitStore,
) *InvoicePaymentService {
	return &InvoicePaymentService{
		db:                  db,
		invoicePaymentStore: invoicePaymentStore,
		transactionStore:    transactionStore,
		accountStore:        accountStore,
		invoiceStore:        invoiceStore,
		splitStore:          splitStore,
	}
}

/*
|---------------------------------------------------------------------------
| Queries
|---------------------------------------------------------------------------
*/

func (s *InvoicePaymentService) GetAll(
	ctx context.Context,
	buildingID int64,
	startDate *string,
	endDate *string,
	peopleID *int,
	status *string,
) ([]dto.InvoicePaymentDto, error) {
	payments, err := s.invoicePaymentStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		return nil, err
	}
	
	var paymentsDto []dto.InvoicePaymentDto
	for _, payment := range payments {
		fmt.Println("payment", payment.AmountCents)
		paymentsDto = append(paymentsDto, dto.InvoicePaymentDto{
			ID: payment.ID,
			TransactionID: payment.TransactionID,
			Reference: payment.Reference,
			Date: payment.Date,
			InvoiceID: payment.InvoiceID,
			UserID: payment.UserID,
			AccountID: payment.AccountID,
			Amount: money.FormatMoneyFromCents(payment.AmountCents),
			Status: payment.Status,
		})
	}

	return paymentsDto, nil
}

func (s *InvoicePaymentService) GetByID(
	ctx context.Context,
	id int64,
) (*dto.InvoicePaymentResponse, error) {
	
	payment, err := s.invoicePaymentStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}


	// initialize paymentDto with empty then check
	var paymentDto dto.InvoicePaymentDto
	if payment != nil {
		paymentDto = dto.MapInvoicePaymentToDto(*payment)
	}

	splits, err := s.splitStore.GetAll(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	// splits to dto
	 splitsDto := dto.MapSplitsToDto(splits)
	

	transaction, err := s.transactionStore.GetByID(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	invoice, err := s.invoiceStore.GetByID(ctx, payment.InvoiceID)
	if err != nil {
		return nil, err
	}

	var invoiceDto dto.InvoiceDto
	if invoice != nil {
		invoiceDto = dto.MapInvoiceToDto(*invoice)
	}

	return &dto.InvoicePaymentResponse{
		Payment: paymentDto,
		Splits: splitsDto,
		Transaction: *transaction,
		Invoice: invoiceDto,
	}, nil
	
}

func (s *InvoicePaymentService) Create(ctx context.Context, paymentDTO dto.CreateInvoicePaymentRequest) (*dto.InvoicePaymentResponse, error) {
	var response dto.InvoicePaymentResponse

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {

		invoice, err := s.invoiceStore.GetByID(ctx, int64(paymentDTO.InvoiceID))
		if err != nil {
			return fmt.Errorf("invoice not found: %v", err)
		}

		// Validate invoice belongs to the building
		if invoice.BuildingID != paymentDTO.BuildingID {
			return fmt.Errorf("invoice does not belong to the specified building")
		}

		arAccount, err := s.accountStore.GetByID(ctx, int64(invoice.ARAccountID))
		if err != nil {
			return fmt.Errorf("A/R account not found: %v", err)
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

		// create transaction
		transaction := &store.Transaction{
			Type:              "payment",
			TransactionDate:   paymentDTO.Date,
			TransactionNumber: paymentDTO.Reference,
			Memo:              paymentDTO.Reference,
			Status:            "1",
			BuildingID:        paymentDTO.BuildingID,
			UnitID:            invoice.UnitID,
			UserID:            1, // TODO: get user id from jwt
		}
		transactionId, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		amountStr := strconv.FormatFloat(paymentDTO.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("invalid amount: %v", err)
		}

		// create splits
		debitSplit := store.Split{
			TransactionID: *transactionId,
			AccountID:     assetAccount.ID,
			Debit:         &paymentDTO.Amount,
			Credit:        nil,
			DebitCents:    &amountCents,
			CreditCents:   nil,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &debitSplit)
		if err != nil {
			return err
		}

		creditSplit := store.Split{
			TransactionID: *transactionId,
			AccountID:     arAccount.ID,
			Credit:        &paymentDTO.Amount,
			Debit:         nil,
			DebitCents:    nil,
			CreditCents:   &amountCents,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &creditSplit)
		if err != nil {
			return err
		}

		// create invoice payment
		invoicePayment := &store.InvoicePayment{
			TransactionID: *transactionId,
			Reference:     paymentDTO.Reference,
			Date:          paymentDTO.Date,
			InvoiceID:     int64(paymentDTO.InvoiceID),
			UserID:        1, // TODO: get user id from jwt
			AccountID:     int64(paymentDTO.AccountID),
			Amount:        paymentDTO.Amount,
			AmountCents:   amountCents,
			Status:        "1",
		}

		_,err = s.invoicePaymentStore.Create(ctx, tx, invoicePayment)
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





func (s *InvoicePaymentService) Update(
	ctx context.Context,
	req dto.UpdateInvoicePaymentRequest,
	paymentID int64,
) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// fetch existing payment
		existing, err := s.invoicePaymentStore.GetByIDTx(ctx, tx, paymentID)
		if err != nil {
			return fmt.Errorf("invoice payment not found: %v", err)
		}

		invoice, err := s.invoiceStore.GetByID(ctx, existing.InvoiceID)
		if err != nil {
			return fmt.Errorf("invoice not found: %v", err)
		}

		// validate account
		if _, err := s.accountStore.GetByID(ctx, int64(req.AccountID)); err != nil {
			return fmt.Errorf("account not found")
		}

		amountStr := strconv.FormatFloat(req.Amount, 'f', -1, 64)
		amountCents, err := money.ParseUSDAmount(amountStr)
		if err != nil {
			return fmt.Errorf("invalid amount: %v", err)
		}

		// update invoice payment
		updatedPayment := &store.InvoicePayment{
			ID:            paymentID,
			TransactionID: existing.TransactionID,
			Reference:     req.Reference,
			Date:          req.Date,
			InvoiceID:     existing.InvoiceID,
			UserID:       1, // TODO: get user id from jwt
			AccountID:     int64(req.AccountID),
			Amount:        req.Amount,
			AmountCents:   amountCents,
			Status:        strconv.Itoa(req.Status),
		}

		if _, err := s.invoicePaymentStore.Update(ctx, tx, updatedPayment); err != nil {
			return err
		}

		// update transaction
		transaction := &store.Transaction{
			ID:                existing.TransactionID,
			Type:              "payment",
			TransactionDate:   req.Date,
			TransactionNumber: req.Reference,
			Memo:              "",
			Status:            strconv.Itoa(req.Status),
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            nil,
		}

		if _, err := s.transactionStore.Update(ctx, tx, transaction); err != nil {
			return err
		}

		// delete existing splits
		err = s.splitStore.DeleteByTransactionID(ctx, tx, existing.TransactionID)
		if err != nil {
			return err
		}

		// create splits
		debitSplit := store.Split{
			TransactionID: existing.TransactionID,
			AccountID:     int64(req.AccountID),
			Debit:         &req.Amount,
			Credit:        nil,
			DebitCents:    &amountCents,
			CreditCents:   nil,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &debitSplit)
		if err != nil {
			return err
		}

		creditSplit := store.Split{
			TransactionID: existing.TransactionID,
			AccountID:     int64(req.AccountID),
			Credit:        &req.Amount,
			Debit:         nil,
			DebitCents:    nil,
			CreditCents:   &amountCents,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			Status:        "1",
		}
		err = s.splitStore.Create(ctx, tx, &creditSplit)
		if err != nil {
			return err
		}

		return nil
	})
}

