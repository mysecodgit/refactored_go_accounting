package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/store"
)

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

type InvoicePaymentDto struct {
	ID int64 `json:"id"`

	TransactionID int64     `json:"transaction_id"`
	Reference     string    `json:"reference"`
	Date          string `json:"date"`

	InvoiceID int64 `json:"invoice_id"`
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`

	Amount string `json:"amount"`
	Status string  `json:"status"` // enum('0','1')
}



type InvoicePaymentResponse struct {
	Payment     InvoicePaymentDto `json:"payment"`
	Splits      []SplitDto        `json:"splits"`
	Transaction store.Transaction    `json:"transaction"`
	Invoice     InvoiceDto        `json:"invoice"`
	ARAccount   *store.Account       `json:"ar_account,omitempty"`
}



func MapInvoicePaymentToDto(p store.InvoicePayment) InvoicePaymentDto {
	return InvoicePaymentDto{
		ID: p.ID,
		TransactionID: p.TransactionID,
		Reference: p.Reference,
		Date: p.Date,
		InvoiceID: p.InvoiceID,
		UserID: p.UserID,
		AccountID: p.AccountID,
		Amount: money.FormatMoneyFromCents(p.AmountCents),
		Status: p.Status,
	}
}

func MapInvoicePaymentsToDto(payments []store.InvoicePayment) []InvoicePaymentDto {
	var dto []InvoicePaymentDto
	for _, p := range payments {
		dto = append(dto, MapInvoicePaymentToDto(p))
	}
	return dto
}