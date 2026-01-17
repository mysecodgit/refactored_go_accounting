package service

import (
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type Service struct {
	User  *UserService
	Building *BuildingService
	Unit *UnitService
	PeopleType *PeopleTypeService
	People *PeopleService
	AccountType *AccountTypeService
	Account *AccountService
	Item *ItemService
	Invoice *InvoiceService
	Reading *ReadingService
	CreditMemo *CreditMemoService
	Check *CheckService
	Journal *JournalService
	InvoicePayment *InvoicePaymentService
	SalesReceipt *SalesReceiptService
	Lease *LeaseService
}

func NewService(
	store store.Storage,
	db *sql.DB,
) *Service {
	return &Service{
		User: NewUserService(store.User),
		Building: NewBuildingService(store.Building),
		Unit: NewUnitService(store.Unit),
		PeopleType: NewPeopleTypeService(store.PeopleType),
		People: NewPeopleService(store.People),
		AccountType: NewAccountTypeService(store.AccountType),
		Account: NewAccountService(store.Account),
		Item: NewItemService(store.Item),
		Invoice: NewInvoiceService(
			db,
			store.Account,
			store.Invoice,
			store.InvoiceItem,
			store.InvoiceAppliedCredit,
			store.InvoiceAppliedDiscount,
			store.InvoicePayment,
			store.Split,
			store.Transaction,
			store.Item,
		),
		Reading: NewReadingService(store.Reading),
		CreditMemo: NewCreditMemoService(
			db,
			store.CreditMemo,
			store.Transaction,
			store.Split,
			store.Account,
		),
		Check: NewCheckService(db, store.Check, store.ExpenseLine, store.Split, store.Transaction, store.Account),
		Journal: NewJournalService(db, store.Journal, store.JournalLine, store.Transaction, store.Split, store.Account),
		InvoicePayment: NewInvoicePaymentService(
			db,
			store.InvoicePayment,
			store.Transaction,
			store.Account,
			store.Invoice,
			store.Split,
		),
		SalesReceipt: NewSalesReceiptService(
			db,
			store.SalesReceipt,
			store.ReceiptItem,
			store.Transaction,
			store.Split,
			store.Account,
			store.Item,
		),
		Lease: NewLeaseService(
			db,
			store.Lease,
			store.Unit,
			store.LeaseFile,
		),
	}
}