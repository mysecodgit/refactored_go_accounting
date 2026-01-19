package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeOutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
)

type Storage struct {
	User *UserStore
	Building *BuildingStore
	Unit *UnitStore
	PeopleType *PeopleTypeStore
	People *PeopleStore
	AccountType *AccountTypeStore
	Account *AccountStore
	Item *ItemStore
	Invoice *InvoiceStore
	InvoiceItem *InvoiceItemStore
	InvoiceAppliedCredit *InvoiceAppliedCreditStore
	InvoiceAppliedDiscount *InvoiceAppliedDiscountStore
	InvoicePayment *InvoicePaymentStore
	Split *SplitStore
	Transaction *TransactionStore
	Reading *ReadingStore
	CreditMemo *CreditMemoStore
	Check *CheckStore
	ExpenseLine *ExpenseLineStore
	Journal *JournalStore
	JournalLine *JournalLineStore
	SalesReceipt *SalesReceiptStore
	ReceiptItem *ReceiptItemStore
	Lease *LeaseStore
	LeaseFile *LeaseFileStore
	Report *ReportStore
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User: &UserStore{db},
		Building: &BuildingStore{db},
		Unit: &UnitStore{db},
		PeopleType: &PeopleTypeStore{db},
		People: &PeopleStore{db},
		AccountType: &AccountTypeStore{db},
		Account: &AccountStore{db},
		Item: &ItemStore{db},
		Invoice: &InvoiceStore{db},
		InvoiceItem: &InvoiceItemStore{db},
		InvoiceAppliedCredit: &InvoiceAppliedCreditStore{db},
		InvoiceAppliedDiscount: &InvoiceAppliedDiscountStore{db},
		InvoicePayment: &InvoicePaymentStore{db},
		Split: &SplitStore{db},
		Transaction: &TransactionStore{db},
		Reading: &ReadingStore{db},
		CreditMemo: &CreditMemoStore{db},
		Check: &CheckStore{db},
		ExpenseLine: &ExpenseLineStore{db},
		Journal: &JournalStore{db},
		JournalLine: &JournalLineStore{db},
		SalesReceipt: &SalesReceiptStore{db},
		ReceiptItem: &ReceiptItemStore{db},
		Lease: &LeaseStore{db},
		LeaseFile: &LeaseFileStore{db},
		Report: &ReportStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}