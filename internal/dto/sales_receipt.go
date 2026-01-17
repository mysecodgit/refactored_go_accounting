package dto

import "github.com/mysecodgit/go_accounting/internal/store"

type ReceiptItemInput struct {
	ItemID        int      `json:"item_id"`
	Qty           *float64 `json:"qty"`
	Rate          *string  `json:"rate"`
	Total         *float64 `json:"total"` // Use manually edited total if provided
	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
}

type SalesReceiptPayload struct {
	ReceiptNo   int             `json:"receipt_no"`
	ReceiptDate string             `json:"receipt_date"`
	UnitID      *int64               `json:"unit_id"`
	PeopleID    *int64               `json:"people_id"`
	AccountID   int64                `json:"account_id"` // Asset account (cash/bank)
	Amount      float64            `json:"amount"`
	Description string             `json:"description"`
	Status      *int               `json:"status"`
	BuildingID  int64                `json:"building_id"`
	Items       []ReceiptItemInput `json:"items"`
}

type CreateSalesReceiptRequest struct {
	SalesReceiptPayload
}

type UpdateSalesReceiptRequest struct {
	ID int `json:"id"`
	SalesReceiptPayload
}

type SalesReceiptResponse struct {
	Receipt     store.SalesReceipt           `json:"receipt"`
	Items       []store.ReceiptItem `json:"items"`
	Splits      []store.Split              `json:"splits"`
	Transaction store.Transaction    `json:"transaction"`
}
