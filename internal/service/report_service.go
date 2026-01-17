package service

import (
	"context"
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type ReportsService struct {
	accountStore    store.AccountStore
	splitStore       store.SplitStore
	transactionStore store.TransactionStore
	invoiceStore    store.InvoiceStore
	paymentStore    store.InvoicePaymentStore
	peopleStore     store.PeopleStore
	peopleTypeStore  store.PeopleTypeStore
	db              *sql.DB
}

func NewReportsService(
	accountStore store.AccountStore,
	splitStore store.SplitStore,
	transactionStore store.TransactionStore,
	invoiceStore store.InvoiceStore,
	paymentStore store.InvoicePaymentStore,
	peopleStore store.PeopleStore,
	peopleTypeStore store.PeopleTypeStore,
	db *sql.DB,
) *ReportsService {
	return &ReportsService{
		accountStore:     accountStore,
		splitStore:       splitStore,
		transactionStore: transactionStore,
		invoiceStore:     invoiceStore,
		paymentStore:     paymentStore,
		peopleStore:      peopleStore,
		peopleTypeStore:  peopleTypeStore,
		db:              db,
	}
}

// GetBalanceSheet generates a balance sheet report
func (s *ReportsService) GetBalanceSheet(ctx context.Context, req dto.BalanceSheetRequest) (*dto.BalanceSheetResponse, error) {
	return nil, nil
}

