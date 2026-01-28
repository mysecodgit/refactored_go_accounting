package store

import (
	"context"
	"database/sql"
	"strings"
)

type ReportStore struct {
	db *sql.DB
}

func NewReportStore(db *sql.DB) *ReportStore {
	return &ReportStore{db: db}
}

type AccountBalance struct {
	AccountID         int     `json:"account_id"`
	AccountNumber     string  `json:"account_number"`
	AccountName       string  `json:"account_name"`
	AccountType       string  `json:"account_type"`
	AccountTypeStatus string  `json:"account_type_status"`
	Debit             float64 `json:"debit"`
	Credit            float64 `json:"credit"`
	Balance           float64 `json:"balance"`
}

func (s *ReportStore) GetBalanceSheet(ctx context.Context, buildingID int, asOfDate string) ([]AccountBalance, error) {
	query := `
		SELECT
    a.id AS account_id,
    a.account_number,
    a.account_name,
    at.type,
    at.typeStatus,
    COALESCE(SUM(s.debit), 0)  AS total_debit,
    COALESCE(SUM(s.credit), 0) AS total_credit,
    CASE
        WHEN LOWER(at.typeStatus) = 'debit'
            THEN COALESCE(SUM(s.debit), 0) - COALESCE(SUM(s.credit), 0)
        ELSE
            COALESCE(SUM(s.credit), 0) - COALESCE(SUM(s.debit), 0)
    END AS balance
FROM accounts a
JOIN account_types at ON a.account_type = at.id
LEFT JOIN (
    SELECT s.*
    FROM splits s
    JOIN transactions t
        ON s.transaction_id = t.id
       AND s.status = '1'
       AND t.transaction_date <= ?
) s ON s.account_id = a.id
WHERE a.building_id = ?
GROUP BY
    a.id,
    a.account_number,
    a.account_name,
    at.type,
    at.typeStatus
ORDER BY a.account_number;



	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, asOfDate, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountBalances []AccountBalance
	for rows.Next() {
		var accountBalance AccountBalance

		if err := rows.Scan(
			&accountBalance.AccountID,
			&accountBalance.AccountNumber,
			&accountBalance.AccountName,
			&accountBalance.AccountType,
			&accountBalance.AccountTypeStatus,
			&accountBalance.Debit,
			&accountBalance.Credit,
			&accountBalance.Balance,
		); err != nil {
			return nil, err
		}
		accountBalances = append(accountBalances, accountBalance)
	}

	return accountBalances, nil
}

type TrialBalanceAccount struct {
	AccountID         int     `json:"account_id"`
	AccountNumber     int     `json:"account_number"`
	AccountName       string  `json:"account_name"`
	AccountType       string  `json:"account_type"`
	AccountTypeStatus string  `json:"account_type_status"`
	DebitBalance      float64 `json:"debit_balance"`          // Debit balance (0 if credit account)
	CreditBalance     float64 `json:"credit_balance"`         // Credit balance (0 if debit account)
	IsTotalRow        bool    `json:"is_total_row,omitempty"` // Flag to indicate this is a total row
}

func (s *ReportStore) GetTrialBalance(ctx context.Context, buildingID int, asOfDate string) ([]TrialBalanceAccount, error) {
	query := `
		SELECT
    a.id AS account_id,
    a.account_number,
    a.account_name,
    at.type AS account_type,
	at.typeName,
    -- Populate total_debit and total_credit based on balance
    CASE 
        WHEN LOWER(at.typeStatus) = 'debit' AND (COALESCE(SUM(s.debit),0) - COALESCE(SUM(s.credit),0)) >= 0 THEN
            COALESCE(SUM(s.debit),0) - COALESCE(SUM(s.credit),0)
        WHEN LOWER(at.typeStatus) = 'credit' AND (COALESCE(SUM(s.credit),0) - COALESCE(SUM(s.debit),0)) < 0 THEN
            -(COALESCE(SUM(s.credit),0) - COALESCE(SUM(s.debit),0))
        ELSE 0
    END AS total_debit,

    CASE 
        WHEN LOWER(at.typeStatus) = 'credit' AND (COALESCE(SUM(s.credit),0) - COALESCE(SUM(s.debit),0)) >= 0 THEN
            COALESCE(SUM(s.credit),0) - COALESCE(SUM(s.debit),0)
        WHEN LOWER(at.typeStatus) = 'debit' AND (COALESCE(SUM(s.debit),0) - COALESCE(SUM(s.credit),0)) < 0 THEN
            -(COALESCE(SUM(s.debit),0) - COALESCE(SUM(s.credit),0))
        ELSE 0
    END AS total_credit

FROM accounts a
JOIN account_types as at ON a.account_type = at.id

LEFT JOIN (
    SELECT s.*
    FROM splits s
    JOIN transactions t
      ON s.transaction_id = t.id
     AND s.status = '1'
     AND t.status = '1'
     AND t.transaction_date <= ?
) s ON s.account_id = a.id

WHERE a.building_id = ?

GROUP BY
    a.id,
    a.account_number,
    a.account_name,
    at.type

ORDER BY a.account_number;

	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, asOfDate, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trialBalanceAccounts []TrialBalanceAccount
	for rows.Next() {
		var trialBalanceAccount TrialBalanceAccount
		if err := rows.Scan(
			&trialBalanceAccount.AccountID,
			&trialBalanceAccount.AccountNumber,
			&trialBalanceAccount.AccountName,
			&trialBalanceAccount.AccountType,
			&trialBalanceAccount.AccountTypeStatus,
			&trialBalanceAccount.DebitBalance,
			&trialBalanceAccount.CreditBalance,
		); err != nil {
			return nil, err
		}
		trialBalanceAccounts = append(trialBalanceAccounts, trialBalanceAccount)
	}
	return trialBalanceAccounts, nil
}

type CustomerSummary struct {
	PeopleID   int     `json:"people_id"`
	PeopleName string  `json:"people_name"`
	Balance    float64 `json:"balance"`
}

func (s *ReportStore) GetCustomerBalanceSummary(ctx context.Context, buildingID int, asOfDate string) ([]CustomerSummary, error) {
	query := `
		SELECT p.id,p.name,(ifnull(SUM(s.debit),0) - ifnull(SUM(s.credit),0)) as balance
FROM people p
LEFT JOIN splits s ON s.people_id = p.id and s.status = "1"
LEFT JOIN transactions t on s.transaction_id = t.id and t.status = "1" 
LEFT JOIN accounts a on s.account_id = a.id
LEFT JOIN account_types as at on a.account_type = at.id
WHERE at.typeName = "Account Receivable" and t.transaction_date <= ? and p.building_id = ?
GROUP BY p.id
HAVING balance <> 0
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, asOfDate, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerSummaries []CustomerSummary
	for rows.Next() {
		var customerSummary CustomerSummary
		if err := rows.Scan(
			&customerSummary.PeopleID,
			&customerSummary.PeopleName,
			&customerSummary.Balance,
		); err != nil {
			return nil, err
		}
		customerSummaries = append(customerSummaries, customerSummary)
	}
	return customerSummaries, nil
}

// get customer balance details by account
type CustomerBalanceDetail struct {
	PeopleID          int
	Name              string
	AccountID         int
	AccountNumber     int
	AccountName       string
	TransactionDate   string
	TransactionNumber string
	Type              string
	Memo              string
	Debit             *float64
	Credit            *float64
}

func (s *ReportStore) GetCustomerBalanceDetail(ctx context.Context, buildingID int, asOfDate string, peopleID *int) ([]CustomerBalanceDetail, error) {
	query := `
		SELECT p.id people_id ,p.name,ac.id account_id,ac.account_number,ac.account_name,t.transaction_date,t.transaction_number,t.type,t.memo,s.debit,s.credit FROM splits s
LEFT JOIN transactions t on s.transaction_id = t.id
LEFT JOIN accounts ac on s.account_id = ac.id
LEFT JOIN account_types as at on ac.account_type = at.id
LEFT JOIN people p on s.people_id = p.id
WHERE s.status = '1' and t.transaction_date <= ? and t.building_id = ? and at.typeName = "Account Receivable"														
	`

	args := []any{asOfDate, buildingID}

	if peopleID != nil {
		query += " AND p.id = ?"
		args = append(args, *peopleID)
	}

	query += " ORDER BY p.id,t.transaction_date asc"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerBalanceDetails []CustomerBalanceDetail
	for rows.Next() {
		var customerBalanceDetail CustomerBalanceDetail
		if err := rows.Scan(
			&customerBalanceDetail.PeopleID,
			&customerBalanceDetail.Name,
			&customerBalanceDetail.AccountID,
			&customerBalanceDetail.AccountNumber,
			&customerBalanceDetail.AccountName,
			&customerBalanceDetail.TransactionDate,
			&customerBalanceDetail.TransactionNumber,
			&customerBalanceDetail.Type,
			&customerBalanceDetail.Memo,
			&customerBalanceDetail.Debit,
			&customerBalanceDetail.Credit,
		); err != nil {
			return nil, err
		}
		customerBalanceDetails = append(customerBalanceDetails, customerBalanceDetail)
	}
	return customerBalanceDetails, nil
}

type TransactionDetail struct {
	PeopleID          *int
	Name              *string
	AccountID         int
	AccountNumber     int
	AccountName       string
	AccountType       string
	TransactionDate   string
	TransactionNumber string
	Type              string
	Memo              string
	Debit             *float64
	Credit            *float64
}

func (s *ReportStore) GetTransactionDetails(ctx context.Context, buildingID int, startDate string, endDate string, accountID []int, unitID *int) ([]TransactionDetail, error) {
	query := `
 SELECT p.id people_id ,p.name,ac.id account_id,ac.account_number,ac.account_name,at.typeName account_type,t.transaction_date,t.transaction_number,t.type,t.memo,s.debit,s.credit FROM splits s
LEFT JOIN transactions t on s.transaction_id = t.id
LEFT JOIN accounts ac on s.account_id = ac.id
LEFT JOIN account_types as at on ac.account_type = at.id
LEFT JOIN people p on s.people_id = p.id
WHERE s.status = '1' and t.transaction_date between ? and ? and t.building_id = 1 
`

	args := []any{startDate, endDate}

	if len(accountID) > 0 {
		// 1. Create one ? for each accountID
		placeholders := make([]string, len(accountID))
		argsSlice := make([]any, len(accountID))
		for i, id := range accountID {
			placeholders[i] = "?"
			argsSlice[i] = id
		}

		// 2. Add to query
		query += " AND ac.id IN (" + strings.Join(placeholders, ",") + ")"

		// 3. Append values to args
		args = append(args, argsSlice...)
	}

	if unitID != nil {
		query += " AND s.unit_id = ?"
		args = append(args, *unitID)
	}

	query += " ORDER BY t.transaction_date,ac.account_number asc"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionDetails []TransactionDetail
	for rows.Next() {
		var transactionDetail TransactionDetail
		if err := rows.Scan(
			&transactionDetail.PeopleID,
			&transactionDetail.Name,
			&transactionDetail.AccountID,
			&transactionDetail.AccountNumber,
			&transactionDetail.AccountName,
			&transactionDetail.AccountType,
			&transactionDetail.TransactionDate,
			&transactionDetail.TransactionNumber,
			&transactionDetail.Type,
			&transactionDetail.Memo,
			&transactionDetail.Debit,
			&transactionDetail.Credit,
		); err != nil {
			return nil, err
		}
		transactionDetails = append(transactionDetails, transactionDetail)
	}
	return transactionDetails, nil
}

type PLAccountRow struct {
	AccountNumber int     `json:"account_number"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"typeName"`
	TotalDebit    float64 `json:"total_debit"`
	TotalCredit   float64 `json:"total_credit"`
	Balance       float64 `json:"balance"`
}

func (s *ReportStore) GetAccountBalanceByAccountType(ctx context.Context, buildingID int, startDate string, endDate string, accountType string) ([]PLAccountRow, error) {
	query := `
		SELECT 
	ac.account_number,
    ac.account_name,
    at.typeName,
    SUM(IFNULL(s.debit, 0)) AS total_debit,
    SUM(IFNULL(s.credit, 0)) AS total_credit,
    SUM(
        CASE 
            WHEN at.typeStatus = 'debit' THEN IFNULL(s.debit, 0) - IFNULL(s.credit, 0)
            WHEN at.typeStatus = 'credit' THEN IFNULL(s.credit, 0) - IFNULL(s.debit, 0)
            ELSE 0
        END
    ) AS balance
FROM splits s
JOIN accounts ac ON s.account_id = ac.id
JOIN account_types at ON ac.account_type = at.id
JOIN transactions t ON s.transaction_id = t.id
WHERE s.status = '1'
  AND t.status = '1'
  AND t.transaction_date BETWEEN ? AND ?
  AND t.building_id = ?
  AND at.type = ?
GROUP BY ac.id, ac.account_name, ac.account_type, at.typeStatus
ORDER BY ac.account_number;
	`

	args := []any{startDate, endDate, buildingID, accountType}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plAccounts []PLAccountRow
	for rows.Next() {
		var plAccount PLAccountRow
		if err := rows.Scan(
			&plAccount.AccountNumber,
			&plAccount.AccountName,
			&plAccount.AccountType,
			&plAccount.TotalDebit,
			&plAccount.TotalCredit,
			&plAccount.Balance,
		); err != nil {
			return nil, err
		}
		plAccounts = append(plAccounts, plAccount)
	}
	return plAccounts, nil
}
