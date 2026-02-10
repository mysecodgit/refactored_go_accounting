package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"

	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type InvoiceStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]store.InvoiceSummary, error)
	GetByID(ctx context.Context, id int64) (*store.Invoice, error)
	Create(ctx context.Context, tx *sql.Tx, invoice *store.Invoice) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, invoice *store.Invoice) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

type InvoiceItemStore interface {
	GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]store.InvoiceItem, error)
	GetByID(ctx context.Context, id int64) (*store.InvoiceItem, error)
	Create(ctx context.Context, tx *sql.Tx, invoiceItem *store.InvoiceItem) error
	Update(ctx context.Context, invoiceItem *store.InvoiceItem) error
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
	DeleteByInvoiceID(ctx context.Context, tx *sql.Tx, invoiceID int64) error
}

type InvoiceAppliedCreditStore interface {
	GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]store.InvoiceAppliedCredit, error)
	GetAllByCreditMemoID(ctx context.Context, creditMemoID int64) ([]store.InvoiceAppliedCredit, error)
	GetByID(ctx context.Context, id int64) (*store.InvoiceAppliedCredit, error)
	Create(ctx context.Context, invoiceAppliedCredit *store.InvoiceAppliedCredit) error
	Update(ctx context.Context, invoiceAppliedCredit *store.InvoiceAppliedCredit) error
	Delete(ctx context.Context, id int64) error
}

type InvoiceAppliedDiscountStore interface {
	GetAllByInvoiceID(ctx context.Context, invoiceID int64) ([]store.InvoiceAppliedDiscount, error)
	GetByID(ctx context.Context, id int64) (*store.InvoiceAppliedDiscount, error)
	Create(ctx context.Context, tx *sql.Tx, invoiceAppliedDiscount *store.InvoiceAppliedDiscount) error
	Update(ctx context.Context, invoiceAppliedDiscount *store.InvoiceAppliedDiscount) error
	Delete(ctx context.Context, id int64) error
}

type TransactionStore interface {
	GetByID(ctx context.Context, id int64) (*store.Transaction, error)
	Create(ctx context.Context, tx *sql.Tx, transaction *store.Transaction) (*int64, error)
	Update(ctx context.Context, tx *sql.Tx, transaction *store.Transaction) (*int64, error)
	Delete(ctx context.Context, id int64) error
}

type SplitStore interface {
	GetAll(ctx context.Context, transactionID int64) ([]store.Split, error)
	GetByTransactionID(ctx context.Context, transactionID int64) ([]store.Split, error)
	GetByID(ctx context.Context, id int64) (*store.Split, error)
	Create(ctx context.Context, tx *sql.Tx, split *store.Split) error
	Update(ctx context.Context, split *store.Split) error
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
	DeleteByTransactionID(ctx context.Context, tx *sql.Tx, transactionID int64) error
}

type InvoiceService struct {
	db                          *sql.DB
	creditMemoStore             CreditMemoStore
	accountStore                AccountStore
	invoiceStore                InvoiceStore
	invoiceItemStore            InvoiceItemStore
	invoiceAppliedCreditStore   InvoiceAppliedCreditStore
	invoiceAppliedDiscountStore InvoiceAppliedDiscountStore
	invoicePaymentStore         InvoicePaymentStore
	splitStore                  SplitStore
	transactionStore            TransactionStore
	itemStore                   ItemStore
}

func NewInvoiceService(
	db *sql.DB,
	creditMemoStore CreditMemoStore,
	accountStore AccountStore,
	invoiceStore InvoiceStore,
	invoiceItemStore InvoiceItemStore,
	invoiceAppliedCreditStore InvoiceAppliedCreditStore,
	invoiceAppliedDiscountStore InvoiceAppliedDiscountStore,
	invoicePaymentStore InvoicePaymentStore,
	splitStore SplitStore,
	transactionStore TransactionStore,
	itemStore ItemStore,
) *InvoiceService {
	return &InvoiceService{
		db:                          db,
		creditMemoStore:             creditMemoStore,
		accountStore:                accountStore,
		invoiceStore:                invoiceStore,
		invoiceItemStore:            invoiceItemStore,
		invoiceAppliedCreditStore:   invoiceAppliedCreditStore,
		invoiceAppliedDiscountStore: invoiceAppliedDiscountStore,
		invoicePaymentStore:         invoicePaymentStore,
		splitStore:                  splitStore,
		transactionStore:            transactionStore,
		itemStore:                   itemStore,
	}
}

func (s *InvoiceService) GetAll(ctx context.Context, buildingID int64, startDate, endDate *string, peopleID *int, status *string) ([]dto.InvoiceListResponse, error) {
	invoices, err := s.invoiceStore.GetAll(ctx, buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		return nil, err
	}
	var invoiceListResponses []dto.InvoiceListResponse
	for _, invoice := range invoices {
		invoiceListResponses = append(invoiceListResponses, dto.InvoiceListResponse{
			ID:                    invoice.ID,
			InvoiceNo:             invoice.InvoiceNo,
			TransactionID:         invoice.TransactionID,
			SalesDate:             invoice.SalesDate,
			DueDate:               invoice.DueDate,
			ARAccountID:           invoice.ARAccountID,
			UnitID:                invoice.UnitID,
			PeopleID:              invoice.PeopleID,
			UserID:                invoice.UserID,
			Amount:                money.FormatMoneyFromCents(invoice.AmountCents),
			Description:           invoice.Description,
			CancelReason:          invoice.CancelReason,
			Status:                invoice.Status,
			BuildingID:            invoice.BuildingID,
			CreatedAt:             invoice.CreatedAt,
			UpdatedAt:             invoice.UpdatedAt,
			PaidAmount:            money.FormatMoneyFromCents(invoice.PaidAmountCents),
			AppliedCreditsTotal:   invoice.AppliedCreditsTotal,
			AppliedDiscountsTotal: invoice.AppliedDiscountsTotal,
			People:                invoice.People,
			Unit:                  invoice.Unit,
		})
	}
	return invoiceListResponses, nil
}

func (s *InvoiceService) GetByID(ctx context.Context, id int64) (map[string]any, error) {
	invoice, err := s.invoiceStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	invoiceItems, err := s.invoiceItemStore.GetAllByInvoiceID(ctx, id)
	if err != nil {
		return nil, err
	}

	var invoiceItemsDto []dto.InvoiceItemDto
	for _, item := range invoiceItems {
		var previousValue *string = nil
		var currentValue *string = nil
		if item.PreviousValueCents != nil {
			previousValueStr := money.FormatScaled5(*item.PreviousValueCents)
			previousValue = &previousValueStr
		}
		if item.CurrentValueCents != nil {
			currentValueStr := money.FormatScaled5(*item.CurrentValueCents)
			currentValue = &currentValueStr
		}

		fmt.Println("*********************** previousValue", item.PreviousValueCents)
		fmt.Println("*********************** currentValue", item.CurrentValueCents)

		//item qty print
		fmt.Println("*********************** item.QtyScaled", item.QtyScaled)
		fmt.Println("*********************** item.RateScaled", item.RateScaled)
		fmt.Println("*********************** item.TotalCents", item.TotalCents)

		invoiceItemsDto = append(invoiceItemsDto, dto.InvoiceItemDto{
			ID:            item.ID,
			InvoiceID:     item.InvoiceID,
			ItemID:        item.ItemID,
			ItemName:      item.ItemName,
			PreviousValue: previousValue,
			CurrentValue:  currentValue,
			Qty:           money.FormatScaled5(item.QtyScaled),
			Rate:          money.FormatScaled5(item.RateScaled),
			Total:         money.FormatMoneyFromCents(item.TotalCents),
			Status:        item.Status,
		})
	}
	appliedCredits, err := s.invoiceAppliedCreditStore.GetAllByInvoiceID(ctx, id)
	if err != nil {
		return nil, err
	}
	appliedDiscounts, err := s.invoiceAppliedDiscountStore.GetAllByInvoiceID(ctx, id)
	if err != nil {
		return nil, err
	}
	payments, err := s.invoicePaymentStore.GetAllByInvoiceID(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Println("*********************** invoice.TransactionID", invoice.TransactionID)

	splits, err := s.splitStore.GetAll(ctx, invoice.TransactionID)
	if err != nil {
		return nil, err
	}

	var splitsDto []dto.SplitDto
	for _, split := range splits {
		var debit *string = nil
		var credit *string = nil
		if split.DebitCents != nil {
			debitStr := money.FormatMoneyFromCents(*split.DebitCents)
			debit = &debitStr
		}
		if split.CreditCents != nil {
			creditStr := money.FormatMoneyFromCents(*split.CreditCents)
			credit = &creditStr
		}

		splitsDto = append(splitsDto, dto.SplitDto{
			ID:            split.ID,
			TransactionID: split.TransactionID,
			AccountID:     split.AccountID,
			Debit:         debit,
			Credit:        credit,
			UnitID:        split.UnitID,
			PeopleID:      split.PeopleID,
			Status:        split.Status,
			CreatedAt:     split.CreatedAt,
			UpdatedAt:     split.UpdatedAt,
			Account:       split.Account,
			Unit:          split.Unit,
			People:        split.People,
		})
	}
	return map[string]any{
		"invoice": dto.InvoiceDto{
			ID:            invoice.ID,
			InvoiceNo:     invoice.InvoiceNo,
			TransactionID: invoice.TransactionID,
			SalesDate:     invoice.SalesDate,
			DueDate:       invoice.DueDate,
			ARAccountID:   invoice.ARAccountID,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			UserID:        invoice.UserID,
			Amount:        money.FormatMoneyFromCents(invoice.AmountCents),
			Description:   invoice.Description,
			CancelReason:  invoice.CancelReason,
			Status:        invoice.Status,
			BuildingID:    invoice.BuildingID,
			CreatedAt:     invoice.CreatedAt,
			UpdatedAt:     invoice.UpdatedAt,
			ARAccount:     invoice.ARAccount,
			Unit:          invoice.Unit,
			People:        invoice.People,
		},
		"items":            invoiceItemsDto,
		"appliedCredits":   appliedCredits,
		"appliedDiscounts": appliedDiscounts,
		"payments":         payments,
		"splits":           splitsDto,
	}, nil
}

func (s *InvoiceService) GetPayments(ctx context.Context, id int64) ([]store.InvoicePayment, error) {
	return s.invoicePaymentStore.GetAllByInvoiceID(ctx, id)
}

func (s *InvoiceService) PreviewInvoiceSplits(ctx context.Context, invoiceDTO dto.InvoicePayloadDTO) ([]store.Split, error) {
	splits, err := s.GenerateInvoiceSplits(ctx, invoiceDTO)
	if err != nil {
		return nil, err
	}

	if err := validateBalanced(splits); err != nil {
		return nil, err
	}
	return splits, nil
}

func (s *InvoiceService) Create(ctx context.Context, invoiceDTO dto.CreateInvoiceRequestDTO) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create transaction
		transaction := &store.Transaction{
			Type:              "invoice",
			TransactionDate:   invoiceDTO.SalesDate,
			TransactionNumber: invoiceDTO.InvoiceNo,
			Memo:              invoiceDTO.Description,
			Status:            "1",
			BuildingID:        invoiceDTO.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            &invoiceDTO.UnitID,
		}
		transactionId, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			fmt.Println("*********************** error creating transaction", err)
			return err
		}

		// create splits

		splits, err := s.GenerateInvoiceSplits(ctx, invoiceDTO.InvoicePayloadDTO)
		if err != nil {
			fmt.Println("*********************** error generating splits", err)
			return err
		}

		if err := validateBalanced(splits); err != nil {
			fmt.Println("*********************** error validating splits", err)
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionId
			err := s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				fmt.Println("*********************** error creating splits", err)
				return err
			}
		}

		amountCents, err := money.ParseUSDAmount(strconv.FormatFloat(invoiceDTO.Amount, 'f', -1, 64))
		if err != nil {
			fmt.Println("*********************** error parsing amount", err)
			return err
		}

		// create invoice
		invoice := &store.Invoice{
			TransactionID: *transactionId,
			InvoiceNo:     invoiceDTO.InvoiceNo,
			SalesDate:     invoiceDTO.SalesDate,
			DueDate:       invoiceDTO.DueDate,
			UnitID:        &invoiceDTO.UnitID,
			PeopleID:      &invoiceDTO.PeopleID,
			ARAccountID:   invoiceDTO.ARAccountID,
			Amount:        invoiceDTO.Amount,
			AmountCents:   amountCents, // TODO : make the amount string on request
			Description:   invoiceDTO.Description,
			Status:        invoiceDTO.Status,
			BuildingID:    invoiceDTO.BuildingID,
			UserID:        1, // TODO: get user id from jwt
		}

		invoiceId, err := s.invoiceStore.Create(ctx, tx, invoice)
		if err != nil {
			fmt.Println("*********************** error converting line input", err)
			return err
		}

		// create invoice items
		for _, item := range invoiceDTO.Items {
			itemrow, err := s.itemStore.GetByID(ctx, int64(item.ItemID))
			if err != nil {
				return err
			}

			var previousValue *string = nil
			var currentValue *string = nil
			if item.PreviousValue != nil {
				previousValueStr := strconv.FormatFloat(*item.PreviousValue, 'f', -1, 64)
				previousValue = &previousValueStr
			}
			if item.CurrentValue != nil {
				currentValueStr := strconv.FormatFloat(*item.CurrentValue, 'f', -1, 64)
				currentValue = &currentValueStr
			}

			lineResult, err := money.ConvertLineInput(money.LineInput{
				Qty:           strconv.FormatFloat(item.Qty, 'f', -1, 64),
				Rate:          strconv.FormatFloat(item.Rate, 'f', -1, 64),
				PreviousValue: previousValue,
				CurrentValue:  currentValue,
			})

			if err != nil {
				fmt.Println("*********************** error converting line input", err)
				return err
			}

			invoiceItem := &store.InvoiceItem{
				InvoiceID:          *invoiceId,
				ItemID:             item.ItemID,
				Qty:                item.Qty,
				Rate:               item.Rate,
				Total:              item.Total,
				PreviousValue:      item.PreviousValue,
				CurrentValue:       item.CurrentValue,
				ItemName:           itemrow.Name,
				QtyScaled:          lineResult.QtyScaled,
				RateScaled:         lineResult.RateScaled,
				TotalCents:         lineResult.TotalCents,
				PreviousValueCents: lineResult.PreviousValueScaled,
				CurrentValueCents:  lineResult.CurrentValueScaled,
			}

			err = s.invoiceItemStore.Create(ctx, tx, invoiceItem)
			if err != nil {
				fmt.Println("*********************** error creating invoice item", err)
				return err
			}
		}

		return nil

	})
}

func (s *InvoiceService) Update(ctx context.Context, invoiceDTO dto.UpdateInvoiceRequestDTO) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// get existing invoice
		existingInvoice, err := s.invoiceStore.GetByID(ctx, int64(invoiceDTO.ID))
		if err != nil {
			fmt.Println("*********************** error getting existing invoice", err)
			return err
		}

		// create transaction
		transaction := &store.Transaction{
			ID:                existingInvoice.TransactionID,
			Type:              "invoice",
			TransactionDate:   invoiceDTO.SalesDate,
			TransactionNumber: invoiceDTO.InvoiceNo,
			Memo:              invoiceDTO.Description,
			Status:            "1",
			BuildingID:        invoiceDTO.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            &invoiceDTO.UnitID,
		}
		transactionId, err := s.transactionStore.Update(ctx, tx, transaction)
		if err != nil {
			fmt.Println("*********************** error creating transaction", err)
			return err
		}

		fmt.Println("*********************** deleting splits --- TransactionID : ", existingInvoice.TransactionID)
		// delete existing splits
		err = s.splitStore.DeleteByTransactionID(ctx, tx, existingInvoice.TransactionID)
		if err != nil {
			fmt.Println("*********************** error deleting splits", err)
			return err
		}

		// create splits

		splits, err := s.GenerateInvoiceSplits(ctx, invoiceDTO.InvoicePayloadDTO)
		if err != nil {
			fmt.Println("*********************** error generating splits", err)
			return err
		}

		if err := validateBalanced(splits); err != nil {
			fmt.Println("*********************** error validating splits", err)
			return err
		}

		for _, split := range splits {
			split.TransactionID = *transactionId
			err := s.splitStore.Create(ctx, tx, &split)
			if err != nil {
				fmt.Println("*********************** error creating splits", err)
				return err
			}
		}

		amountCents, err := money.ParseUSDAmount(strconv.FormatFloat(invoiceDTO.Amount, 'f', -1, 64))
		if err != nil {
			fmt.Println("*********************** error parsing amount", err)
			return err
		}

		// update invoice
		invoice := &store.Invoice{
			ID:            existingInvoice.ID,
			TransactionID: *transactionId,
			InvoiceNo:     invoiceDTO.InvoiceNo,
			SalesDate:     invoiceDTO.SalesDate,
			DueDate:       invoiceDTO.DueDate,
			UnitID:        &invoiceDTO.UnitID,
			PeopleID:      &invoiceDTO.PeopleID,
			ARAccountID:   invoiceDTO.ARAccountID,
			Amount:        invoiceDTO.Amount,
			AmountCents:   amountCents, // TODO : make the amount string on request
			Description:   invoiceDTO.Description,
			BuildingID:    invoiceDTO.BuildingID,
			UserID:        1, // TODO: get user id from jwt
		}

		invoiceId, err := s.invoiceStore.Update(ctx, tx, invoice)
		if err != nil {
			fmt.Println("*********************** error updating invoice", err)
			return err
		}

		// delete existing invoice items
		err = s.invoiceItemStore.DeleteByInvoiceID(ctx, tx, existingInvoice.ID)
		if err != nil {
			fmt.Println("*********************** error deleting invoice items", err)
			return err
		}

		// create invoice items
		for _, item := range invoiceDTO.Items {
			itemrow, err := s.itemStore.GetByID(ctx, int64(item.ItemID))
			if err != nil {
				fmt.Println("*********************** error getting item", err)
				return err
			}

			var previousValue *string = nil
			var currentValue *string = nil
			if item.PreviousValue != nil {
				previousValueStr := strconv.FormatFloat(*item.PreviousValue, 'f', -1, 64)
				previousValue = &previousValueStr
			}
			if item.CurrentValue != nil {
				currentValueStr := strconv.FormatFloat(*item.CurrentValue, 'f', -1, 64)
				currentValue = &currentValueStr
			}

			lineResult, err := money.ConvertLineInput(money.LineInput{
				Qty:           strconv.FormatFloat(item.Qty, 'f', -1, 64),
				Rate:          strconv.FormatFloat(item.Rate, 'f', -1, 64),
				PreviousValue: previousValue,
				CurrentValue:  currentValue,
			})

			invoiceItem := &store.InvoiceItem{
				InvoiceID:          *invoiceId,
				ItemID:             item.ItemID,
				Qty:                item.Qty,
				Rate:               item.Rate,
				Total:              item.Total,
				PreviousValue:      item.PreviousValue,
				CurrentValue:       item.CurrentValue,
				ItemName:           itemrow.Name,
				QtyScaled:          lineResult.QtyScaled,
				RateScaled:         lineResult.RateScaled,
				TotalCents:         lineResult.TotalCents,
				PreviousValueCents: lineResult.PreviousValueScaled,
				CurrentValueCents:  lineResult.CurrentValueScaled,
			}

			err = s.invoiceItemStore.Create(ctx, tx, invoiceItem)
			if err != nil {
				fmt.Println("*********************** error creating invoice item", err)
				return err
			}
		}

		return nil

	})
}

type splitAccumulator struct {
	Debit       float64
	Credit      float64
	DebitCents  int64
	CreditCents int64
}

func (s *InvoiceService) GenerateInvoiceSplits(
	ctx context.Context,
	req dto.InvoicePayloadDTO,
) ([]store.Split, error) {

	acc := make(map[int64]*splitAccumulator)

	// Helper function
	addDebit := func(accountID int64, amount float64, amountCents int64) {
		if acc[accountID] == nil {
			acc[accountID] = &splitAccumulator{}
		}
		acc[accountID].Debit += amount
		acc[accountID].DebitCents += amountCents
	}

	addCredit := func(accountID int64, amount float64, amountCents int64) {
		if acc[accountID] == nil {
			acc[accountID] = &splitAccumulator{}
		}
		acc[accountID].Credit += amount
		acc[accountID].CreditCents += amountCents
	}

	// 1. AR Debit
	var amountCents int64 = 0

	// 2. Item lines
	for _, line := range req.Items {

		item, err := s.itemStore.GetByID(ctx, int64(line.ItemID))
		if err != nil {
			return nil, err
		}

		lineResult, err := money.ConvertLineInput(money.LineInput{Qty: strconv.FormatFloat(line.Qty, 'f', -1, 64), Rate: strconv.FormatFloat(line.Rate, 'f', -1, 64)})
		if err != nil {
			return nil, err
		}

		amountCents += lineResult.TotalCents

		lineTotal := line.Total
		if lineTotal == 0 {
			lineTotal = line.Qty * line.Rate
		}

		totalCents := lineResult.TotalCents

		switch item.Type {

		case "service":
			addCredit(*item.IncomeAccount, lineTotal, totalCents)

		case "discount":
			addDebit(*item.IncomeAccount, lineTotal, totalCents)

		case "payment":
			// reduces AR via asset account
			addCredit(*item.AssetAccount, lineTotal, totalCents)

		default:
			return nil, fmt.Errorf("unsupported item type: %s", item.Type)
		}
	}

	addDebit(int64(req.ARAccountID), req.Amount, amountCents)
	// 3. Build splits
	splits := make([]store.Split, 0, len(acc))

	for accountID, v := range acc {

		var debit, credit *float64
		var debitCents, creditCents *int64
		if v.Debit > 0 {
			d := v.Debit
			debit = &d
			debitCents = &v.DebitCents
		}
		if v.Credit > 0 {
			c := v.Credit
			credit = &c
			creditCents = &v.CreditCents
		}

		splits = append(splits, store.Split{
			AccountID:   accountID,
			Debit:       debit,
			Credit:      credit,
			DebitCents:  debitCents,
			CreditCents: creditCents,
			UnitID:      &req.UnitID,
			PeopleID:    &req.PeopleID,
			Status:      "1",
		})
	}

	return splits, nil
}

func validateBalanced(splits []store.Split) error {
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

	if math.Abs(debit-credit) > 0.0001 {
		return fmt.Errorf("unbalanced entry: debit %.2f ≠ credit %.2f", debit, credit)
	}

	if debitCents != creditCents {
		return fmt.Errorf("unbalanced entry: debit cents %d ≠ credit cents %d", debitCents, creditCents)
	}

	return nil
}

// PAYMENT RELATED FUNCTIONS

func (s *InvoiceService) CreateInvoicePayment(ctx context.Context, paymentDTO dto.CreateInvoicePaymentRequest) (*dto.InvoicePaymentResponse, error) {
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

		// create splits
		debitSplit := store.Split{
			TransactionID: *transactionId,
			AccountID:     assetAccount.ID,
			Debit:         &paymentDTO.Amount,
			Credit:        nil,
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
			Status:        "1",
		}

		_, err = s.invoicePaymentStore.Create(ctx, tx, invoicePayment)
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

func (s *InvoiceService) GetInvoiceDiscounts(ctx context.Context, id int64) ([]store.InvoiceAppliedDiscount, error) {
	return s.invoiceAppliedDiscountStore.GetAllByInvoiceID(ctx, id)
}

func (s *InvoiceService) CreateInvoiceDiscount(ctx context.Context, invoiceID int64, discountDTO dto.CreateInvoiceAppliedDiscountRequest) (*dto.InvoiceAppliedDiscountResponse, error) {

	var response dto.InvoiceAppliedDiscountResponse

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {

		// Get invoice
		invoice, err := s.invoiceStore.GetByID(ctx, invoiceID)

		fmt.Println("*********************** invoice", invoiceID)
		if err != nil {
			return fmt.Errorf("invoice not found: %v", err)
		}

		if invoice.PeopleID == nil {
			return fmt.Errorf("invoice must have a people_id")
		}

		// Validate A/R account matches
		if invoice.ARAccountID != discountDTO.ARAccount {
			return fmt.Errorf("A/R account does not match invoice A/R account")
		}

		// Validate accounts exist
		_, err = s.accountStore.GetByID(ctx, int64(discountDTO.ARAccount))
		if err != nil {
			return fmt.Errorf("A/R account not found: %v", err)
		}

		_, err = s.accountStore.GetByID(ctx, int64(discountDTO.IncomeAccount))
		if err != nil {
			return fmt.Errorf("income account not found: %v", err)
		}

		transactionMemo := fmt.Sprintf("Discount applied to Invoice #%s: %s", invoice.InvoiceNo, discountDTO.Description)

		// create transaction
		transaction := &store.Transaction{
			Type:              "payment", // TODO: change to discount later or something else
			TransactionDate:   discountDTO.Date,
			TransactionNumber: discountDTO.Reference,
			Memo:              transactionMemo,
			Status:            "1",
			BuildingID:        invoice.BuildingID,
			UserID:            1, // TODO: get user id from jwt
			UnitID:            invoice.UnitID,
		}

		transactionId, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// create splits
		debitSplit := store.Split{
			TransactionID: *transactionId,
			AccountID:     int64(discountDTO.ARAccount),
			Credit:        &discountDTO.Amount,
			Debit:         nil,
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
			AccountID:     int64(discountDTO.IncomeAccount),
			Debit:         &discountDTO.Amount,
			Credit:        nil,
			UnitID:        invoice.UnitID,
			PeopleID:      invoice.PeopleID,
			Status:        "1",
		}

		err = s.splitStore.Create(ctx, tx, &creditSplit)
		if err != nil {
			return err
		}

		// create invoice applied discount
		invoiceAppliedDiscount := &store.InvoiceAppliedDiscount{
			Reference:       discountDTO.Reference,
			TransactionID:   *transactionId,
			InvoiceID:       invoiceID,
			Amount:          discountDTO.Amount,
			Description:     discountDTO.Description,
			Date:            discountDTO.Date,
			Status:          "1",
			ARAccountID:     int64(discountDTO.ARAccount),
			IncomeAccountID: int64(discountDTO.IncomeAccount),
		}

		err = s.invoiceAppliedDiscountStore.Create(ctx, tx, invoiceAppliedDiscount)
		if err != nil {
			return err
		}

		response.InvoiceAppliedDiscount = *invoiceAppliedDiscount
		response.Splits = []store.Split{debitSplit, creditSplit}
		response.Transaction = *transaction

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (s *InvoiceService) GetAppliedCredits(ctx context.Context, invoiceID int64) ([]store.InvoiceAppliedCredit, error) {
	return s.invoiceAppliedCreditStore.GetAllByInvoiceID(ctx, invoiceID)
}

// GET INVOICE AVAILABLE CREDITS
func (s *InvoiceService) GetInvoiceAvailableCredits(ctx context.Context, invoiceID int64) (*dto.AvailableCreditsResponse, error) {

	// get invoice
	invoice, err := s.invoiceStore.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %v", err)
	}

	if invoice.PeopleID == nil {
		return &dto.AvailableCreditsResponse{
			InvoiceID: int(invoiceID),
			PeopleID:  0,
			Credits:   []dto.AvailableCreditMemo{},
		}, nil
	}

	// get available credits
	credits, err := s.creditMemoStore.GetByPeopleID(ctx, *invoice.PeopleID)
	if err != nil {
		return nil, fmt.Errorf("available credits not found: %v", err)
	}

	peopleID := *invoice.PeopleID
	availableCredits := []dto.AvailableCreditMemo{}

	for _, credit := range credits {
		if credit.PeopleID == peopleID && strconv.Itoa(credit.Status) == "1" { // TODO: change status to string later
			// Get applied amount for this credit memo
			appliedCredits, err := s.invoiceAppliedCreditStore.GetAllByCreditMemoID(ctx, credit.ID)
			if err != nil {
				continue
			}

			// sum applied credits amount
			appliedAmount := 0.0
			for _, appliedCredit := range appliedCredits {
				appliedAmount += appliedCredit.Amount
			}

			// Calculate available amount and round to 2 decimal places to avoid floating-point precision issues
			availableAmount := credit.Amount - appliedAmount
			// Round to 2 decimal places
			availableAmount = float64(int(availableAmount*100+0.5)) / 100

			if availableAmount > 0 {
				availableCredits = append(availableCredits, dto.AvailableCreditMemo{
					ID:              int(credit.ID),
					Date:            credit.Date,
					Amount:          credit.Amount,
					AppliedAmount:   appliedAmount,
					AvailableAmount: availableAmount,
					Description:     credit.Description,
				})
			}
		}
	}

	return &dto.AvailableCreditsResponse{
		InvoiceID: int(invoiceID),
		PeopleID:  int(*invoice.PeopleID),
		Credits:   availableCredits,
	}, nil
}

// APPLY INVOICE CREDITS
func (s *InvoiceService) ApplyInvoiceCredits(ctx context.Context, req dto.CreateInvoiceAppliedCreditRequest) error {
	// Validate amount
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Get invoice
	invoice, err := s.invoiceStore.GetByID(ctx, int64(req.InvoiceID))
	if err != nil {
		return fmt.Errorf("invoice not found: %v", err)
	}

	if invoice.PeopleID == nil {
		return fmt.Errorf("invoice must have a people_id")
	}

	// if invoice.ARAccountID == nil {
	// 	return fmt.Errorf("invoice must have an A/R account")
	// }

	// Get credit memo
	creditMemo, err := s.creditMemoStore.GetByID(ctx, int64(req.CreditMemoID))
	if err != nil {
		return fmt.Errorf("credit memo not found: %v", err)
	}

	// Validate people_id matches
	if creditMemo.PeopleID != *invoice.PeopleID {
		return fmt.Errorf("credit memo people_id does not match invoice people_id")
	}

	// Check available amount
	appliedCredits, err := s.invoiceAppliedCreditStore.GetAllByInvoiceID(ctx, int64(req.InvoiceID))

	// sum applied credits amount
	appliedAmount := 0.0
	for _, appliedCredit := range appliedCredits {
		appliedAmount += appliedCredit.Amount
	}

	if err != nil {
		return fmt.Errorf("failed to get applied amount: %v", err)
	}

	availableAmount := creditMemo.Amount - appliedAmount
	if req.Amount > availableAmount {
		return fmt.Errorf("amount exceeds available credit. Available: %.2f, Requested: %.2f", availableAmount, req.Amount)
	}

	// Create invoice applied credit record (no transaction or splits needed)
	appliedCreditStatus := "1"
	appliedCredit := store.InvoiceAppliedCredit{
		InvoiceID:    int64(req.InvoiceID),
		CreditMemoID: int64(req.CreditMemoID),
		Amount:       req.Amount,
		Description:  req.Description,
		Date:         req.Date,
		Status:       appliedCreditStatus,
	}

	err = s.invoiceAppliedCreditStore.Create(ctx, &appliedCredit)
	if err != nil {
		return fmt.Errorf("failed to create invoice applied credit: %v", err)
	}

	return nil
}
