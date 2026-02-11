package dto

import money "github.com/mysecodgit/go_accounting/internal/accounting"
import store "github.com/mysecodgit/go_accounting/internal/store"

type InvoiceAppliedDiscountDto struct {
	ID int64 `json:"id"`

	Reference     string `json:"reference"`
	InvoiceID     int64  `json:"invoice_id"`
	TransactionID int64  `json:"transaction_id"`

	ARAccountID     int64 `json:"ar_account"`
	IncomeAccountID int64 `json:"income_account"`

	Amount      string `json:"amount"`
	Description string `json:"description"`
	Date        string `json:"date"`

	Status string `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// map invoice applied discount to dto
func MapInvoiceAppliedDiscountToDto(invoiceAppliedDiscount store.InvoiceAppliedDiscount) *InvoiceAppliedDiscountDto {
	return &InvoiceAppliedDiscountDto{
		ID:              invoiceAppliedDiscount.ID,
		Reference:       invoiceAppliedDiscount.Reference,
		InvoiceID:       invoiceAppliedDiscount.InvoiceID,
		TransactionID:   invoiceAppliedDiscount.TransactionID,
		ARAccountID:     invoiceAppliedDiscount.ARAccountID,
		IncomeAccountID: invoiceAppliedDiscount.IncomeAccountID,
		Amount:          money.FormatMoneyFromCents(invoiceAppliedDiscount.AmountCents),
		Description:     invoiceAppliedDiscount.Description,
		Date:            invoiceAppliedDiscount.Date,
		Status:          invoiceAppliedDiscount.Status,
		CreatedAt:       invoiceAppliedDiscount.CreatedAt,
		UpdatedAt:       invoiceAppliedDiscount.UpdatedAt,
	}
}

// map invoice applied discounts to dto
func MapInvoiceAppliedDiscountsToDto(invoiceAppliedDiscounts []store.InvoiceAppliedDiscount) []*InvoiceAppliedDiscountDto {
	var dto []*InvoiceAppliedDiscountDto
	for _, invoiceAppliedDiscount := range invoiceAppliedDiscounts {
		dto = append(dto, MapInvoiceAppliedDiscountToDto(invoiceAppliedDiscount))
	}
	return dto
}