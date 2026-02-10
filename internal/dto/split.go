package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	store "github.com/mysecodgit/go_accounting/internal/store"
)

type SplitDto struct {
	ID            int64    `json:"id"`
	TransactionID int64    `json:"transaction_id"`
	AccountID     int64    `json:"account_id"`
	Debit         *string `json:"debit"`
	Credit        *string `json:"credit"`
	UnitID        *int64   `json:"unit_id"`
	PeopleID      *int64   `json:"people_id"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	

	// relationships
	Account store.Account `json:"account"`
	Unit    store.Unit    `json:"unit"`
	People  store.People  `json:"people"`
}


// map split to dto
func MapSplitToDto(s store.Split) SplitDto {
	var debitFormatted *string
	var creditFormatted *string
	if s.DebitCents != nil {
		debit := money.FormatMoneyFromCents(*s.DebitCents)
		debitFormatted = &debit
	}
	if s.CreditCents != nil {
		credit := money.FormatMoneyFromCents(*s.CreditCents)
		creditFormatted = &credit
	}

	return SplitDto{
		ID:            s.ID,
		TransactionID: s.TransactionID,
		AccountID:     s.AccountID,
		Debit:         debitFormatted,
		Credit:        creditFormatted,
		UnitID:        s.UnitID,
		PeopleID:      s.PeopleID,
		Status:        s.Status,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

// map splits to dto
func MapSplitsToDto(splits []store.Split) []SplitDto {
	var dto []SplitDto
	for _, s := range splits {
		dto = append(dto, MapSplitToDto(s))
	}
	return dto
}