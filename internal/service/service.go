package service

import (
	"database/sql"

	"github.com/mysecodgit/go_accounting/internal/store"
)

type Service struct {
	Auth             *AuthService
	User             *UserService
	Building         *BuildingService
	Unit             *UnitService
	PeopleType       *PeopleTypeService
	People           *PeopleService
	AccountType      *AccountTypeService
	Account          *AccountService
	Item             *ItemService
	Invoice          *InvoiceService
	Reading          *ReadingService
	CreditMemo       *CreditMemoService
	Check            *CheckService
	Bill             *BillService
	BillPayment      *BillPaymentService
	Journal          *JournalService
	InvoicePayment   *InvoicePaymentService
	SalesReceipt     *SalesReceiptService
	Lease            *LeaseService
	Report           *ReportService
	UserBuilding     *UserBuildingService
	Permission       *PermissionService
	Role             *RoleService
	RolePermission   *RolePermissionService
	UserBuildingRole *UserBuildingRoleService
}

func NewService(
	store store.Storage,
	db *sql.DB,
	jwtSecret string,
) *Service {
	return &Service{
		Auth:        NewAuthService(store.User, jwtSecret),
		User:        NewUserService(store.User),
		Building:    NewBuildingService(store.Building),
		Unit:        NewUnitService(store.Unit),
		PeopleType:  NewPeopleTypeService(store.PeopleType),
		People:      NewPeopleService(store.People),
		AccountType: NewAccountTypeService(store.AccountType),
		Account:     NewAccountService(store.Account),
		Item:        NewItemService(store.Item),
		Invoice: NewInvoiceService(
			db,
			store.CreditMemo,
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
		Reading: NewReadingService(store.Reading, db),
		CreditMemo: NewCreditMemoService(
			db,
			store.CreditMemo,
			store.Transaction,
			store.Split,
			store.Account,
		),
		Check:       NewCheckService(db, store.Check, store.ExpenseLine, store.Split, store.Transaction, store.Account),
		Bill:        NewBillService(db, store.Bill, store.BillExpenseLine, store.Split, store.Transaction, store.Account),
		BillPayment: NewBillPaymentService(db, store.BillPayment, store.Transaction, store.Account, store.Bill, store.Split),
		Journal:     NewJournalService(db, store.Journal, store.JournalLine, store.Transaction, store.Split, store.Account),
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
			store.People,
		),
		Report: NewReportService(
			store.Report,
			store.Unit,
		),
		UserBuilding:     NewUserBuildingService(store.UserBuilding),
		Permission:       NewPermissionService(store.Permission),
		Role:             NewRoleService(store.Role),
		RolePermission:   NewRolePermissionService(store.RolePermission),
		UserBuildingRole: NewUserBuildingRoleService(store.UserBuildingRole),
	}
}
