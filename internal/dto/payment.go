package dto

import "github.com/mysecodgit/go_accounting/internal/store"

type InvoicePaymentPayload struct{
	Reference  string  `json:"reference"`
	Date       string  `json:"date"`
	InvoiceID  int     `json:"invoice_id"`
	AccountID  int     `json:"account_id"` // Asset account (cash/bank)
	Amount     float64 `json:"amount"`
	Status     int    `json:"status"`
	BuildingID int64     `json:"building_id"`
}

type CreateInvoicePaymentRequest struct {
	InvoicePaymentPayload
}

type UpdateInvoicePaymentRequest struct {
	ID int64 `json:"id"`
	InvoicePaymentPayload
}

type InvoicePaymentResponse struct {
	Payment     store.InvoicePayment `json:"payment"`
	Splits      []store.Split        `json:"splits"`
	Transaction store.Transaction    `json:"transaction"`
	Invoice     store.Invoice        `json:"invoice"`
	ARAccount   *store.Account       `json:"ar_account,omitempty"`
}


