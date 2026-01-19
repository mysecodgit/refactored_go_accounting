package dto

import "github.com/mysecodgit/go_accounting/internal/store"

type LeaseWithPeople struct {
	Lease  store.Lease   `json:"lease"`
	People *store.People `json:"people,omitempty"`
}

type ReadingPayload struct {
	ItemID        int      `json:"item_id"`
	UnitID        int      `json:"unit_id"`
	LeaseID       *int64   `json:"lease_id"`
	ReadingMonth  *string  `json:"reading_month"`
	ReadingYear   *string  `json:"reading_year"`
	ReadingDate   string   `json:"reading_date"`
	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
	UnitPrice     *float64 `json:"unit_price"`
	TotalAmount   *float64 `json:"total_amount"`
	Notes         *string  `json:"notes"`
	Status        string   `json:"status"`
}

type CreateReadingRequest struct {
	Readings []ReadingPayload `json:"readings"`
}

type UpdateReadingRequest struct {
	ID int `json:"id"`
	ReadingPayload
}

type ReadingResponse struct {
	Reading store.Reading `json:"reading"`
	Item    store.Item    `json:"item"`
	Unit    store.Unit    `json:"unit"`
	Lease   *store.Lease  `json:"lease,omitempty"`
}

type ReadingListItem struct {
	Reading store.Reading `json:"reading"`
	Item    store.Item    `json:"item"`
	Unit    store.Unit    `json:"unit"`
	Lease   *store.Lease  `json:"lease,omitempty"`
}

type BulkImportReadingRequest struct {
	ItemID        int      `json:"item_id"`
	UnitID        int      `json:"unit_id"`
	LeaseID       *int     `json:"lease_id"`
	ReadingMonth  *string  `json:"reading_month"`
	ReadingYear   *string  `json:"reading_year"`
	ReadingDate   string   `json:"reading_date"`
	PreviousValue *float64 `json:"previous_value"`
	CurrentValue  *float64 `json:"current_value"`
	UnitPrice     *float64 `json:"unit_price"`
	TotalAmount   *float64 `json:"total_amount"`
	Notes         *string  `json:"notes"`
	Status        string   `json:"status"`
}

type BulkImportReadingsRequest struct {
	Readings []BulkImportReadingRequest `json:"readings"`
}

type BulkImportReadingsResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"`
}
