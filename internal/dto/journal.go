package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type JournalDto struct {
	ID            int64      `json:"id"`
	TransactionID int64      `json:"transaction_id"`
	Reference     string     `json:"reference"`
	JournalDate   string  `json:"journal_date"`
	BuildingID    int64      `json:"building_id"`
	Memo          *string    `json:"memo,omitempty"`
	TotalAmount   string   `json:"total_amount,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type JournalLineDto struct {
	ID          int64    `json:"id"`
	JournalID   int64    `json:"journal_id"`
	AccountID   int64    `json:"account_id"`
	UnitID      *int64   `json:"unit_id,omitempty"`
	PeopleID    *int64   `json:"people_id,omitempty"`
	Description *string  `json:"description,omitempty"`
	Debit       string  `json:"debit,omitempty"`
	Credit      string  `json:"credit,omitempty"`
}

type JournalLineInput struct {
	AccountID   int      `json:"account_id"`
	UnitID      *int64     `json:"unit_id"`
	PeopleID    *int64     `json:"people_id"`
	Description *string  `json:"description"`
	Debit       *float64 `json:"debit"`
	Credit      *float64 `json:"credit"`
}

type JournalPayloadDTO struct {
	Reference   string             `json:"reference"`
	JournalDate string             `json:"journal_date"`
	BuildingID  int64                `json:"building_id"`
	Memo        *string            `json:"memo"`
	TotalAmount float64            `json:"total_amount"`
	Lines       []JournalLineInput `json:"lines"`
}

type CreateJournalRequest struct {
	JournalPayloadDTO
}

type UpdateJournalRequest struct {
	ID int `json:"id"`
	JournalPayloadDTO
}

type JournalResponse struct {
	Journal     JournalDto            `json:"journal"`
	Lines       []JournalLineDto      `json:"lines"`
	Splits      []SplitDto               `json:"splits"`
	Transaction store.Transaction `json:"transaction"`
}

type JournalResponseDetails struct {
	Journal      JournalDto       `json:"journal"`
	Lines []*JournalLineDto `json:"lines"`
	Transaction  store.Transaction   `json:"transaction"`
	Splits  []SplitDto   `json:"splits"`
}


type AvailableCreditMemo struct {
	ID              int     `json:"id"`
	Date            string  `json:"date"`
	Amount          string `json:"amount"`
	AppliedAmount   string `json:"applied_amount"`   // Amount already applied to other invoices
	AvailableAmount string `json:"available_amount"` // Amount available to apply
	Description     string  `json:"description"`
}

type AvailableCreditsResponse struct {
	InvoiceID int                  `json:"invoice_id"`
	PeopleID  int                  `json:"people_id"`
	Credits   []AvailableCreditMemo `json:"credits"`
}

type CreateInvoiceAppliedCreditRequest struct {
	InvoiceID    int     `json:"invoice_id"`
	CreditMemoID int     `json:"credit_memo_id"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Date         string  `json:"date"`
}


// map Journal to JournalDto
func MapJournalToJournalDto(j store.Journal) *JournalDto {
	return &JournalDto{
		ID:            j.ID,
		TransactionID: j.TransactionID,
		Reference:     j.Reference,
		JournalDate:   j.JournalDate,
		BuildingID:    j.BuildingID,
		Memo:          j.Memo,
		TotalAmount:   money.FormatMoneyFromCents(j.AmountCents),
		CreatedAt:     j.CreatedAt,
	}
}

// map JournalLine to JournalLineDto
func MapJournalLineToJournalLineDto(l store.JournalLine) *JournalLineDto {
	return &JournalLineDto{
		ID:            l.ID,
		JournalID:     l.JournalID,
		AccountID:     l.AccountID,
		UnitID:        l.UnitID,
		PeopleID:      l.PeopleID,
		Description:   l.Description,
		Debit:         money.FormatMoneyFromCents(l.DebitCents),
		Credit:        money.FormatMoneyFromCents(l.CreditCents),
	}
}

// map []journal to []JournalDto
func MapJournalsToJournalDtos(journals []store.Journal) []*JournalDto {
	var dto []*JournalDto
	for _, j := range journals {
		dto = append(dto, MapJournalToJournalDto(j))
	}
	return dto
}

// map []journal_line to []JournalLineDto
func MapJournalLinesToJournalLineDtos(lines []store.JournalLine) []*JournalLineDto {
	var dto []*JournalLineDto
	for _, l := range lines {
		dto = append(dto, MapJournalLineToJournalLineDto(l))
	}
	return dto
}