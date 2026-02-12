package dto

import (
	money "github.com/mysecodgit/go_accounting/internal/accounting"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type BillExpenseLineInput struct {
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id"`
	PeopleID    *int64  `json:"people_id"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
}

type BillPayloadDTO struct {
	BillNo       string                 `json:"bill_no"`
	BillDate     string                 `json:"bill_date"`
	DueDate      string                 `json:"due_date"`
	APAccountID  int64                  `json:"ap_account_id"`
	UnitID       *int64                 `json:"unit_id"`
	PeopleID     *int64                 `json:"people_id"`
	BuildingID   int64                  `json:"building_id"`
	Amount       float64                `json:"amount"`
	Description  string                 `json:"description"`
	ExpenseLines []BillExpenseLineInput `json:"expense_lines"`
}

type CreateBillRequest struct {
	BillPayloadDTO
}

type UpdateBillRequest struct {
	ID int `json:"id"`
	BillPayloadDTO
}

// Bill Payment DTOs
type BillPaymentPayload struct {
	Reference  string  `json:"reference"`
	Date       string  `json:"date"`
	BillID     int     `json:"bill_id"`
	AccountID  int     `json:"account_id"` // Asset account (cash/bank)
	Amount     float64 `json:"amount"`
	Status     int     `json:"status"`
	BuildingID int64   `json:"building_id"`
}

type CreateBillPaymentRequest struct {
	BillPaymentPayload
}

type UpdateBillPaymentRequest struct {
	ID int64 `json:"id"`
	BillPaymentPayload
}

type BillPaymentResponse struct {
	Payment     BillPaymentDto `json:"payment"`
	Splits      []SplitDto     `json:"splits"`
	Transaction store.Transaction `json:"transaction"`
	Bill        BillDto        `json:"bill"`
	APAccount   *store.Account    `json:"ap_account,omitempty"`
}

type BillDto struct {
	ID            int64   `json:"id"`
	BillNo        string  `json:"bill_no"`
	TransactionID int64   `json:"transaction_id"`
	BillDate      string  `json:"bill_date"`
	DueDate       string  `json:"due_date"`
	APAccountID   int64   `json:"ap_account_id"`
	UnitID        *int64  `json:"unit_id"`
	PeopleID      *int64  `json:"people_id"`
	UserID        int64   `json:"user_id"`
	Amount        string `json:"amount"`
	Description   string  `json:"description"`
	CancelReason  *string `json:"cancel_reason"`
	Status        string  `json:"status"` // enum('0','1')
	BuildingID    int64   `json:"building_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// map store.Bill to BillDto

func MapBillToDto(b store.Bill) *BillDto {
	return &BillDto{
		ID:            b.ID,
		BillNo:        b.BillNo,
		TransactionID: b.TransactionID,
		BillDate:      b.BillDate,
		DueDate:       b.DueDate,
		APAccountID:   b.APAccountID,
		UnitID:        b.UnitID,
		PeopleID:      b.PeopleID,
		UserID:        b.UserID,
		Amount:        money.FormatMoneyFromCents(b.AmountCents),
		Description:   b.Description,
		CancelReason:  b.CancelReason,
		Status:        b.Status,
		BuildingID:    b.BuildingID,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

// map []store.Bill to []BillDto
func MapBillsToDtos(bills []store.Bill) []*BillDto {
	
	var dtoBills []*BillDto
	for _, b := range bills {
		dtoBills = append(dtoBills, MapBillToDto(b))
	}
	return dtoBills
}

type BillExpenseLineDto struct {
	ID          int64   `json:"id"`
	BillID      int64   `json:"bill_id"`
	AccountID   int64   `json:"account_id"`
	UnitID      *int64  `json:"unit_id"`
	PeopleID    *int64  `json:"people_id"`
	Description *string `json:"description"`
	Amount      string `json:"amount"`
}

// map store.BillExpenseLine to BillExpenseLineDto
func MapBillExpenseLineToDto(l store.BillExpenseLine) *BillExpenseLineDto {
	return &BillExpenseLineDto{
		ID:          l.ID,
		BillID:      l.BillID,
		AccountID:   l.AccountID,
		UnitID:      l.UnitID,
		PeopleID:    l.PeopleID,
		Description: l.Description,
		Amount:      money.FormatMoneyFromCents(l.AmountCents),
	}
}

// map []store.BillExpenseLine to []BillExpenseLineDto
func MapBillExpenseLinesToDtos(lines []store.BillExpenseLine) []*BillExpenseLineDto {
	var dtoLines []*BillExpenseLineDto
	for _, l := range lines {
		dtoLines = append(dtoLines, MapBillExpenseLineToDto(l))
	}
	return dtoLines
}


type BillResponseDetails struct {
	Bill        BillDto              `json:"bill"`
	ExpenseLines []*BillExpenseLineDto `json:"expense_lines"`
	Splits       []SplitDto           `json:"splits"`
	Transaction  store.Transaction       `json:"transaction"`
}

type BillPaymentDto struct {
	ID int64 `json:"id"`

	TransactionID int64  `json:"transaction_id"`
	Reference     string `json:"reference"`
	Date          string `json:"date"`

	BillID    int64 `json:"bill_id"`
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`

	Amount string `json:"amount"`
	Status string  `json:"status"` // enum('0','1')

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// map store.BillPayment to BillPaymentDto
func MapBillPaymentToDto(p store.BillPayment) *BillPaymentDto {
	return &BillPaymentDto{
		ID: p.ID,
		TransactionID: p.TransactionID,
		Reference: p.Reference,
		Date: p.Date,
		BillID: p.BillID,
		UserID: p.UserID,
		AccountID: p.AccountID,
		Amount: money.FormatMoneyFromCents(p.AmountCents),
		Status: p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}


// map []store.BillPayment to []BillPaymentDto
func MapBillPaymentsToDtos(payments []store.BillPayment) []*BillPaymentDto {
	var dtoPayments []*BillPaymentDto
	for _, p := range payments {
		dtoPayments = append(dtoPayments, MapBillPaymentToDto(p))
	}
	return dtoPayments
}