package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	store "github.com/mysecodgit/go_accounting/internal/store"
)

type InvoiceAppliedCreditDto struct {
	ID int64 `json:"id"`

	InvoiceID    int64 `json:"invoice_id"`
	CreditMemoID int64 `json:"credit_memo_id"`

	Amount      string `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`

	Status string `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// map invoice applied credit to dto

func MapInvoiceAppliedCreditToDto(invoiceAppliedCredit store.InvoiceAppliedCredit) *InvoiceAppliedCreditDto {
	return &InvoiceAppliedCreditDto{
		ID: invoiceAppliedCredit.ID,
		InvoiceID: invoiceAppliedCredit.InvoiceID,
		CreditMemoID: invoiceAppliedCredit.CreditMemoID,
		Amount: money.FormatMoneyFromCents(invoiceAppliedCredit.AmountCents),
		Description: invoiceAppliedCredit.Description,
		Date: invoiceAppliedCredit.Date,
		Status: invoiceAppliedCredit.Status,
		CreatedAt: invoiceAppliedCredit.CreatedAt,
		UpdatedAt: invoiceAppliedCredit.UpdatedAt,
	}
}

// map invoice applied credits to dto
func MapInvoiceAppliedCreditsToDto(invoiceAppliedCredits []store.InvoiceAppliedCredit) []*InvoiceAppliedCreditDto {
	var dto []*InvoiceAppliedCreditDto
	for _, invoiceAppliedCredit := range invoiceAppliedCredits {
		dto = append(dto, MapInvoiceAppliedCreditToDto(invoiceAppliedCredit))
	}
	return dto
}