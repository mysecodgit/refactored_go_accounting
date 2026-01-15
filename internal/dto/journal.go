package dto

import "github.com/mysecodgit/go_accounting/internal/store"

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
	Journal     store.Journal            `json:"journal"`
	Lines       []store.JournalLine      `json:"lines"`
	Splits      []store.Split               `json:"splits"`
	Transaction store.Transaction `json:"transaction"`
}
