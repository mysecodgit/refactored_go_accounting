package dto

import money "github.com/mysecodgit/go_accounting/internal/accounting"
import store "github.com/mysecodgit/go_accounting/internal/store"

type ExpenseLineInput struct {
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id"`
	PeopleID    *int64  `json:"people_id"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
}

type CheckPayloadDTO struct {
	CheckDate        string             `json:"check_date"`
	ReferenceNumber  *string            `json:"reference_number"`
	PaymentAccountID int64              `json:"payment_account_id"`
	BuildingID       int64              `json:"building_id"`
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

type CheckDto struct {
	ID               int64   `json:"id"`
	TransactionID    int64   `json:"transaction_id"`
	CheckDate        string  `json:"check_date"`
	ReferenceNumber  string  `json:"reference_number"`
	PaymentAccountID int64   `json:"payment_account_id"`
	BuildingID       int64   `json:"building_id"`
	Memo             *string `json:"memo"`
	TotalAmount      string `json:"total_amount"`
	CreatedAt        string  `json:"created_at"`
}

// map from store.Check to CheckDto
func MapCheckToDto(check store.Check) *CheckDto {
	return &CheckDto{
		ID:               check.ID,
		TransactionID:    check.TransactionID,
		CheckDate:        check.CheckDate,
		ReferenceNumber:  check.ReferenceNumber,
		PaymentAccountID: check.PaymentAccountID,
		BuildingID:       check.BuildingID,
		Memo:             check.Memo,
		TotalAmount:      money.FormatMoneyFromCents(check.AmountCents),
		CreatedAt:        check.CreatedAt,
	}
}

// map from []store.Check to []CheckDto
func MapChecksToDtos(checks []store.Check) []*CheckDto {
	var dtos []*CheckDto
	for _, check := range checks {
		dtos = append(dtos, MapCheckToDto(check))
	}
	return dtos
}

type ExpenseLineDto struct {
	ID          int64   `json:"id"`
	CheckID     int64   `json:"check_id"`
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id,omitempty"`
	PeopleID    *int64  `json:"people_id,omitempty"`
	Description *string `json:"description"`
	Amount      string `json:"amount"`
}

// map from store.ExpenseLine to ExpenseLineDto
func MapExpenseLineToDto(expenseLine store.ExpenseLine) *ExpenseLineDto {
	return &ExpenseLineDto{
		ID:          expenseLine.ID,
		CheckID:     expenseLine.CheckID,
		AccountID:   expenseLine.AccountID,
		UnitID:      expenseLine.UnitID,
		PeopleID:    expenseLine.PeopleID,
		Description: expenseLine.Description,
		Amount:      money.FormatMoneyFromCents(expenseLine.AmountCents),
	}
}

// map from []store.ExpenseLine to []ExpenseLineDto
func MapExpenseLinesToDtos(expenseLines []store.ExpenseLine) []*ExpenseLineDto {
	var dtos []*ExpenseLineDto
	for _, expenseLine := range expenseLines {
		dtos = append(dtos, MapExpenseLineToDto(expenseLine))
	}
	return dtos
}

type CheckResponseDetails struct {
	Check        *CheckDto         `json:"check"`
	ExpenseLines []*ExpenseLineDto `json:"expense_lines"`
	Splits       []SplitDto       `json:"splits"`
	Transaction  store.Transaction   `json:"transaction"`
}