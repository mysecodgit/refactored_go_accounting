package dto

import (
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
