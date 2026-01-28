package dto

import "github.com/mysecodgit/go_accounting/internal/store"

type BillExpenseLineInput struct {
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id"`
	PeopleID    *int64  `json:"people_id"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
}

type BillPayloadDTO struct {
	BillNo      string                `json:"bill_no"`
	BillDate    string                `json:"bill_date"`
	DueDate     string                `json:"due_date"`
	APAccountID int64                 `json:"ap_account_id"`
	UnitID      *int64                `json:"unit_id"`
	PeopleID    *int64                `json:"people_id"`
	BuildingID  int64                 `json:"building_id"`
	Amount      float64               `json:"amount"`
	Description string                `json:"description"`
	ExpenseLines []BillExpenseLineInput `json:"expense_lines"`
}

type CreateBillRequest struct {
	BillPayloadDTO
}

type UpdateBillRequest struct {
	ID int `json:"id"`
	BillPayloadDTO
}

// Bill Payment DTOs
type BillPaymentPayload struct {
	Reference  string  `json:"reference"`
	Date       string  `json:"date"`
	BillID     int     `json:"bill_id"`
	AccountID  int     `json:"account_id"` // Asset account (cash/bank)
	Amount     float64 `json:"amount"`
	Status     int     `json:"status"`
	BuildingID int64   `json:"building_id"`
}

type CreateBillPaymentRequest struct {
	BillPaymentPayload
}

type UpdateBillPaymentRequest struct {
	ID int64 `json:"id"`
	BillPaymentPayload
}

type BillPaymentResponse struct {
	Payment     store.BillPayment `json:"payment"`
	Splits      []store.Split      `json:"splits"`
	Transaction store.Transaction  `json:"transaction"`
	Bill        store.Bill          `json:"bill"`
	APAccount   *store.Account      `json:"ap_account,omitempty"`
}
