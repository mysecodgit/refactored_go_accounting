package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	store "github.com/mysecodgit/go_accounting/internal/store"
)

// payload dto
type CreditMemoPayloadDTO struct {
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int64     `json:"people_id"`
	BuildingID       int64     `json:"building_id"`
	UnitID           int64     `json:"unit_id"`
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
}
type CreateCreditMemoRequest struct {
	CreditMemoPayloadDTO
}

type UpdateCreditMemoRequest struct {
	ID               int     `json:"id"`
	CreditMemoPayloadDTO
}

type CreditMemoSummaryDto struct {
	ID               int     `json:"id"`
	TransactionID    int     `json:"transaction_id"`
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	UserID           int     `json:"user_id"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int     `json:"people_id"`
	BuildingID       int     `json:"building_id"`
	UnitID           int     `json:"unit_id"`
	Amount           string `json:"amount"`
	Description      string  `json:"description"`
	Status           int     `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	People           store.People  `json:"people"`
	Unit             store.Unit    `json:"unit"`
	UsedCredits      float64 `json:"used_credits"`
	Balance          float64 `json:"balance"`
}

// map credit memo summary to dto
func MapCreditMemoSummaryToDto(cm store.CreditMemoSummary) CreditMemoSummaryDto {
	return CreditMemoSummaryDto{
		ID: cm.ID,
		TransactionID: cm.TransactionID,
		Reference: cm.Reference,
		Date: cm.Date,
		UserID: cm.UserID,
		DepositTo: cm.DepositTo,
		LiabilityAccount: cm.LiabilityAccount,
		PeopleID: cm.PeopleID,
		BuildingID: cm.BuildingID,
		UnitID: cm.UnitID,
		Amount: money.FormatMoneyFromCents(cm.AmountCents),
		Description: cm.Description,
		Status: cm.Status,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
		People: cm.People,
		Unit: cm.Unit,
		UsedCredits: cm.UsedCredits,
		Balance: cm.Balance,
	}
}

// map credit memo summaries to dto
func MapCreditMemoSummariesToDto(cm []store.CreditMemoSummary) []CreditMemoSummaryDto {
	var dto []CreditMemoSummaryDto
	for _, cm := range cm {
		dto = append(dto, MapCreditMemoSummaryToDto(cm))
	}
	return dto
}


type CreditMemoDto struct {
	ID               int64   `json:"id"`
	TransactionID    int64   `json:"transaction_id"`
	Reference        string  `json:"reference"`
	Date             string  `json:"date"`
	UserID           int64   `json:"user_id"`
	DepositTo        int     `json:"deposit_to"`
	LiabilityAccount int     `json:"liability_account"`
	PeopleID         int64     `json:"people_id"`
	BuildingID       int64     `json:"building_id"`
	UnitID           int64     `json:"unit_id"`
	Amount           string `json:"amount"`
	Description      string  `json:"description"`
	Status           int     `json:"status"` // enum('0','1')
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`

	// relationships
	People store.People `json:"people"`
	Unit   store.Unit   `json:"unit"`
}

// map credit memo to dto
func MapCreditMemoToDto(cm store.CreditMemo) *CreditMemoDto {
	return &CreditMemoDto{
		ID: cm.ID,
		TransactionID: cm.TransactionID,
		Reference: cm.Reference,
		Date: cm.Date,
		UserID: cm.UserID,
		DepositTo: cm.DepositTo,
		LiabilityAccount: cm.LiabilityAccount,
		PeopleID: cm.PeopleID,
		BuildingID: cm.BuildingID,
		UnitID: cm.UnitID,
		Amount: money.FormatMoneyFromCents(cm.AmountCents),
		Description: cm.Description,
		Status: cm.Status,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
		People: cm.People,
		Unit: cm.Unit,
		
	}
}

// map credit memos to dto
func MapCreditMemosToDto(cm []store.CreditMemo) []*CreditMemoDto {
	var dto []*CreditMemoDto
	for _, cm := range cm {
		dto = append(dto, MapCreditMemoToDto(cm))
	}
	return dto
}

type CreditMemoDetailsResponse struct {
	CreditMemo  *CreditMemoDto              `json:"credit_memo"`
	Splits      []SplitDto          `json:"splits"`
	Transaction *store.Transaction `json:"transaction"`
}
