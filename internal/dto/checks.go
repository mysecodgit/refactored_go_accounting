package dto

type ExpenseLineInput struct {
	AccountID   int64     `json:"account_id"`
	UnitID      *int64    `json:"unit_id"`
	PeopleID    *int64    `json:"people_id"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
}

type CheckPayloadDTO struct {
	CheckDate        string             `json:"check_date"`
	ReferenceNumber  *string            `json:"reference_number"`
	PaymentAccountID int64                `json:"payment_account_id"`
	BuildingID       int64                `json:"building_id"`
	Memo             *string            `json:"memo"`
	TotalAmount      float64            `json:"total_amount"`
	ExpenseLines     []ExpenseLineInput `json:"expense_lines"`
}

type CreateCheckRequest struct {
	CheckPayloadDTO
}

type UpdateCheckRequest struct {
	ID int `json:"id"`
	CheckPayloadDTO
}
