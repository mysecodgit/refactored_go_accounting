package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type ReportStore interface {
	GetBalanceSheet(ctx context.Context, buildingID int, asOfDate string) ([]store.AccountBalance, error)
	GetTrialBalance(ctx context.Context, buildingID int, asOfDate string) ([]store.TrialBalanceAccount, error)
	GetCustomerBalanceSummary(ctx context.Context, buildingID int, asOfDate string) ([]store.CustomerSummary, error)
	GetCustomerBalanceDetail(ctx context.Context, buildingID int, asOfDate string, peopleID *int) ([]store.CustomerBalanceDetail, error)
	GetTransactionDetails(ctx context.Context, buildingID int, startDate string, endDate string, accountID []int, unitID *int) ([]store.TransactionDetail, error)
	GetAccountBalanceByAccountType(ctx context.Context, buildingID int, startDate string, endDate string, accountType string) ([]store.PLAccountRow, error)
	GetAccountBalanceByAccountTypeAndUnit(ctx context.Context, buildingID int, startDate string, endDate string, accountType string) ([]store.PLAccountRowByUnit, error)
}

type ReportService struct {
	reportStore ReportStore
	unitStore   UnitStoreInterface
}

type UnitStoreInterface interface {
	GetAll(ctx context.Context, buildingID int64) ([]store.Unit, error)
}

func NewReportService(
	reportStore ReportStore,
	unitStore UnitStoreInterface,
) *ReportService {
	return &ReportService{
		reportStore: reportStore,
		unitStore:   unitStore,
	}
}

// GetBalanceSheet generates a balance sheet report
func (s *ReportService) GetBalanceSheet(ctx context.Context, buildingID int, asOfDate string) (*dto.BalanceSheetResponse, error) {
	fmt.Println("***************************** asOfDate", asOfDate)
	fmt.Println("***************************** buildingID", buildingID)
	accountBalances, err := s.reportStore.GetBalanceSheet(ctx, buildingID, asOfDate)
	if err != nil {
		return nil, err
	}
	fmt.Println("***************************** accountBalances", accountBalances)
	assets := []dto.AccountBalance{}
	liabilities := []dto.AccountBalance{}
	equity := []dto.AccountBalance{}
	incomeAccounts := []dto.AccountBalance{}
	expenseAccounts := []dto.AccountBalance{}

	for _, account := range accountBalances {
		accountType := account.AccountType
		balance := account.Balance

		// Skip accounts with 0 balance
		if balance == 0 {
			continue
		}

		accountBalance := dto.AccountBalance{
			AccountID:     account.AccountID,
			AccountNumber: account.AccountNumber,
			AccountName:   account.AccountName,
			AccountType:   account.AccountType,
			Balance:       balance,
		}

		// Categorize based on account type field (Asset, Liability, Equity, Income, Expense)
		typeLower := strings.ToLower(accountType)
		fmt.Println("***************************** typeLower", typeLower)
		if typeLower == "asset" {
			assets = append(assets, accountBalance)
		} else if typeLower == "liability" {
			liabilities = append(liabilities, accountBalance)
		} else if typeLower == "equity" {
			equity = append(equity, accountBalance)
		} else if typeLower == "income" {
			incomeAccounts = append(incomeAccounts, accountBalance)
		} else if typeLower == "expense" {
			expenseAccounts = append(expenseAccounts, accountBalance)
		}
	}

	// Calculate Net Income = Total Income - Total Expenses
	totalIncome := 0.0
	for _, income := range incomeAccounts {
		totalIncome += income.Balance
	}

	totalExpenses := 0.0
	for _, expense := range expenseAccounts {
		totalExpenses += expense.Balance
	}

	netIncome := totalIncome - totalExpenses

	// Add Net Income to equity section (only if it's not zero)
	if netIncome != 0 {
		equity = append(equity, dto.AccountBalance{
			AccountID:     0, // 0 indicates this is a calculated value, not an actual account
			AccountNumber: "",
			AccountName:   "Net Income",
			AccountType:   "Net Income",
			Balance:       netIncome,
		})
	}

	// Calculate totals
	totalAssets := 0.0
	for _, asset := range assets {
		totalAssets += asset.Balance
	}

	totalLiabilities := 0.0
	for _, liability := range liabilities {
		totalLiabilities += liability.Balance
	}

	totalEquity := 0.0
	for _, eq := range equity {
		totalEquity += eq.Balance
	}

	totalLiabilitiesAndEquity := totalLiabilities + totalEquity

	// Round to 2 decimals for presentation + stable comparisons (avoid floating point artifacts)
	round2 := func(v float64) float64 {
		return math.Round(v*100) / 100
	}
	for i := range assets {
		assets[i].Balance = round2(assets[i].Balance)
	}
	for i := range liabilities {
		liabilities[i].Balance = round2(liabilities[i].Balance)
	}
	for i := range equity {
		equity[i].Balance = round2(equity[i].Balance)
	}
	totalAssets = round2(totalAssets)
	totalLiabilities = round2(totalLiabilities)
	totalEquity = round2(totalEquity)
	totalLiabilitiesAndEquity = round2(totalLiabilitiesAndEquity)

	// Consider balanced if equal at 2-decimal precision
	isBalanced := totalAssets == totalLiabilitiesAndEquity

	return &dto.BalanceSheetResponse{
		BuildingID:                buildingID,
		AsOfDate:                  asOfDate,
		Assets:                    dto.BalanceSheetSection{SectionName: "Assets", Accounts: assets, Total: totalAssets},
		Liabilities:               dto.BalanceSheetSection{SectionName: "Liabilities", Accounts: liabilities, Total: totalLiabilities},
		Equity:                    dto.BalanceSheetSection{SectionName: "Equity", Accounts: equity, Total: totalEquity},
		TotalAssets:               totalAssets,
		TotalLiabilitiesAndEquity: totalLiabilitiesAndEquity,
		IsBalanced:                isBalanced,
	}, nil

}

func (s *ReportService) GetTrialBalance(ctx context.Context, buildingID int, asOfDate string) (*dto.TrialBalanceResponse, error) {
	fmt.Println("***************************** asOfDate", asOfDate)
	fmt.Println("***************************** buildingID", buildingID)
	accounts, err := s.reportStore.GetTrialBalance(ctx, buildingID, asOfDate)
	if err != nil {
		return nil, err
	}
	fmt.Println("***************************** trialBalanceAccounts", accounts)
	trialBalanceAccounts := []dto.TrialBalanceAccount{}
	totalDebit := 0.0
	totalCredit := 0.0

	for _, account := range accounts {

		// Only include accounts that have debit or credit (active splits)
		if account.DebitBalance > 0 || account.CreditBalance > 0 {
			trialBalanceAccounts = append(trialBalanceAccounts, dto.TrialBalanceAccount{
				AccountID:     account.AccountID,
				AccountNumber: account.AccountNumber,
				AccountName:   account.AccountName,
				AccountType:   account.AccountType,
				DebitBalance:  account.DebitBalance,
				CreditBalance: account.CreditBalance,
			})

			totalDebit += account.DebitBalance
			totalCredit += account.CreditBalance
		}
	}

	// Check if balanced (allowing for small rounding differences)
	isBalanced := (totalDebit-totalCredit) < 0.01 && (totalDebit-totalCredit) > -0.01

	// Add total row
	trialBalanceAccounts = append(trialBalanceAccounts, dto.TrialBalanceAccount{
		AccountID:     0,
		AccountNumber: 0,
		AccountName:   "TOTAL",
		AccountType:   "",
		DebitBalance:  totalDebit,
		CreditBalance: totalCredit,
		IsTotalRow:    true,
	})

	return &dto.TrialBalanceResponse{
		BuildingID:  buildingID,
		AsOfDate:    asOfDate,
		Accounts:    trialBalanceAccounts,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
		IsBalanced:  isBalanced,
	}, nil
}

func (s *ReportService) GetCustomerBalanceSummary(ctx context.Context, buildingID int, asOfDate string) (*dto.CustomerBalanceSummaryResponse, error) {
	customers, err := s.reportStore.GetCustomerBalanceSummary(ctx, buildingID, asOfDate)
	if err != nil {
		return nil, err
	}

	// Calculate total balance
	totalBalance := 0.0
	customersList := []dto.CustomerBalance{}
	for _, customer := range customers {
		customersList = append(customersList, dto.CustomerBalance{
			PeopleID:   customer.PeopleID,
			PeopleName: customer.PeopleName,
			Balance:    customer.Balance,
		})
		totalBalance += customer.Balance
	}

	return &dto.CustomerBalanceSummaryResponse{
		BuildingID:   buildingID,
		AsOfDate:     asOfDate,
		Customers:    customersList,
		TotalBalance: totalBalance,
	}, nil

}

func (s *ReportService) GetCustomerBalanceDetail(ctx context.Context, buildingID int, asOfDate string, peopleID *int) (*map[string]any, error) {
	customerBalanceDetails, err := s.reportStore.GetCustomerBalanceDetail(ctx, buildingID, asOfDate, peopleID)
	if err != nil {
		return nil, err
	}

	customerBalanceDetailsList := []CustomerBalanceDetail{}
	for _, customerBalanceDetail := range customerBalanceDetails {
		customerBalanceDetailsList = append(customerBalanceDetailsList, CustomerBalanceDetail{
			PeopleID:          customerBalanceDetail.PeopleID,
			Name:              customerBalanceDetail.Name,
			AccountID:         customerBalanceDetail.AccountID,
			AccountName:       customerBalanceDetail.AccountName,
			AccountNumber:     customerBalanceDetail.AccountNumber,
			TransactionDate:   customerBalanceDetail.TransactionDate,
			TransactionNumber: customerBalanceDetail.TransactionNumber,
			Type:              customerBalanceDetail.Type,
			Memo:              customerBalanceDetail.Memo,
			Debit:             customerBalanceDetail.Debit,
			Credit:            customerBalanceDetail.Credit,
		})
	}

	response := GroupTransactionsWithGrandTotals(customerBalanceDetailsList)

	return &map[string]any{
		"buildingID":          buildingID,
		"asOfDate":            asOfDate,
		"customers":           response.Customers,
		"grand_total_balance": response.GrandTotalBalance,
		"grand_total_credit":  response.GrandTotalCredit,
		"grand_total_debit":   response.GrandTotalDebit,
	}, nil
}

type CustomerBalanceDetail struct {
	PeopleID          int
	Name              string
	AccountID         int
	AccountName       string
	AccountNumber     int
	TransactionDate   string
	TransactionNumber string
	Type              string
	Memo              string
	Debit             *float64
	Credit            *float64
}

type Customer struct {
	PeopleID     int        `json:"people_id"`
	PeopleName   string     `json:"people_name"`
	Accounts     []*Account `json:"accounts"`
	TotalDebit   float64    `json:"total_debit"`
	TotalCredit  float64    `json:"total_credit"`
	TotalBalance float64    `json:"total_balance"`
	IsHeader     bool       `json:"is_header"`
}

type Account struct {
	AccountID     int     `json:"account_id"`
	AccountName   string  `json:"account_name"`
	AccountNumber int     `json:"account_number"`
	Splits        []Split `json:"splits,omitempty"`
	TotalDebit    float64 `json:"total_debit"`
	TotalCredit   float64 `json:"total_credit"`
	TotalBalance  float64 `json:"total_balance"`
}

type Split struct {
	PeopleID          *int     `json:"people_id"`
	PeopleName        *string  `json:"people_name"`
	TransactionNumber string   `json:"transaction_number"`
	TransactionDate   string   `json:"transaction_date"`
	TransactionType   string   `json:"transaction_type"`
	TransactionMemo   string   `json:"transaction_memo"`
	AccountID         int      `json:"account_id"`
	AccountName       string   `json:"account_name"`
	AccountNumber     int      `json:"account_number"`
	Debit             *float64 `json:"debit"`
	Credit            *float64 `json:"credit"`
	Balance           float64  `json:"balance"`
}

type CustomersResponse struct {
	Customers         []Customer `json:"customers"`
	GrandTotalDebit   float64    `json:"grand_total_debit"`
	GrandTotalCredit  float64    `json:"grand_total_credit"`
	GrandTotalBalance float64    `json:"grand_total_balance"`
}

func GroupTransactionsWithGrandTotals(rows []CustomerBalanceDetail) CustomersResponse {
	customersMap := make(map[int]*Customer)
	accountMap := make(map[int]map[int]*Account)

	var grandDebit, grandCredit, grandBalance float64

	for _, r := range rows {

		// ---- Customer ----
		if _, ok := customersMap[r.PeopleID]; !ok {
			customersMap[r.PeopleID] = &Customer{
				PeopleID:   r.PeopleID,
				PeopleName: r.Name,
				IsHeader:   true,
			}
			accountMap[r.PeopleID] = make(map[int]*Account)
		}
		customer := customersMap[r.PeopleID]

		// ---- Account ----
		if _, ok := accountMap[r.PeopleID][r.AccountID]; !ok {
			account := &Account{
				AccountID:     r.AccountID,
				AccountName:   r.AccountName,
				AccountNumber: r.AccountNumber,
			}
			accountMap[r.PeopleID][r.AccountID] = account
			customer.Accounts = append(customer.Accounts, account)
		}
		account := accountMap[r.PeopleID][r.AccountID]

		// ---- Amounts ----
		var debit, credit float64
		if r.Debit != nil {
			debit = *r.Debit
		}
		if r.Credit != nil {
			credit = *r.Credit
		}

		// ---- Running balance ----
		balance := debit - credit
		if len(account.Splits) > 0 {
			balance += account.Splits[len(account.Splits)-1].Balance
		}

		// ---- Split ----
		account.Splits = append(account.Splits, Split{
			TransactionNumber: r.TransactionNumber,
			TransactionDate:   r.TransactionDate,
			TransactionType:   r.Type,
			TransactionMemo:   r.Memo,
			AccountID:         r.AccountID,
			AccountName:       r.AccountName,
			AccountNumber:     r.AccountNumber,
			Debit:             r.Debit,
			Credit:            r.Credit,
			Balance:           balance,
		})

		// ---- Totals (Account) ----
		account.TotalDebit += debit
		account.TotalCredit += credit
		account.TotalBalance += debit - credit

		// ---- Totals (Customer) ----
		customer.TotalDebit += debit
		customer.TotalCredit += credit
		customer.TotalBalance += debit - credit

		// ---- Grand totals ----
		grandDebit += debit
		grandCredit += credit
		grandBalance += debit - credit
	}

	// ---- Map → Slice ----
	var customers []Customer
	for _, c := range customersMap {
		customers = append(customers, *c)
	}

	return CustomersResponse{
		Customers:         customers,
		GrandTotalDebit:   grandDebit,
		GrandTotalCredit:  grandCredit,
		GrandTotalBalance: grandBalance,
	}
}

// getTransactionDetails
func (s *ReportService) GetTransactionDetails(ctx context.Context, buildingID int, startDate string, endDate string, accountID []int, unitID *int) (*map[string]any, error) {
	transactionDetails, err := s.reportStore.GetTransactionDetails(ctx, buildingID, startDate, endDate, accountID, unitID)
	if err != nil {
		return nil, err
	}
	transactionDetailsList := []TransactionDetail{}
	for _, transactionDetail := range transactionDetails {
		transactionDetailsList = append(transactionDetailsList, TransactionDetail{
			PeopleID:          transactionDetail.PeopleID,
			Name:              transactionDetail.Name,
			AccountID:         transactionDetail.AccountID,
			AccountNumber:     transactionDetail.AccountNumber,
			AccountName:       transactionDetail.AccountName,
			AccountType:       transactionDetail.AccountType,
			TransactionDate:   transactionDetail.TransactionDate,
			TransactionNumber: transactionDetail.TransactionNumber,
			Type:              transactionDetail.Type,
			Memo:              transactionDetail.Memo,
			Debit:             transactionDetail.Debit,
			Credit:            transactionDetail.Credit,
		})
	}

	response := BuildLedgerResponse(transactionDetailsList, buildingID, startDate, endDate)
	return &map[string]any{
		"buildingID":         buildingID,
		"startDate":          startDate,
		"endDate":            endDate,
		"grand_total_debit":  response.GrandTotalDebit,
		"grand_total_credit": response.GrandTotalCredit,
		"accounts":           response.Accounts,
	}, nil
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

type LedgerResponse struct {
	BuildingID       int              `json:"building_id"`
	StartDate        string           `json:"start_date"`
	EndDate          string           `json:"end_date"`
	GrandTotalDebit  float64          `json:"grand_total_debit"`
	GrandTotalCredit float64          `json:"grand_total_credit"`
	Accounts         []*AccountLedger `json:"accounts"`
}

type AccountLedger struct {
	AccountID     int     `json:"account_id"`
	AccountNumber int     `json:"account_number"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	Splits        []Split `json:"splits"`
	TotalDebit    float64 `json:"total_debit"`
	TotalCredit   float64 `json:"total_credit"`
	TotalBalance  float64 `json:"total_balance"`
	IsTotalRow    bool    `json:"is_total_row,omitempty"`
}

func BuildLedgerResponse(
	rows []TransactionDetail,
	buildingID int,
	startDate, endDate string,
) LedgerResponse {

	accountMap := make(map[int]*AccountLedger)
	var accounts []*AccountLedger

	var grandDebit, grandCredit float64

	for _, r := range rows {

		// ---- Account ----
		if _, ok := accountMap[r.AccountID]; !ok {
			account := &AccountLedger{
				AccountID:     r.AccountID,
				AccountNumber: r.AccountNumber,
				AccountName:   r.AccountName,
				AccountType:   r.AccountType,
				Splits:        []Split{},
			}
			accountMap[r.AccountID] = account
			accounts = append(accounts, account)
		}
		account := accountMap[r.AccountID]

		// ---- Amounts ----
		var debit, credit float64
		if r.Debit != nil {
			debit = *r.Debit
		}
		if r.Credit != nil {
			credit = *r.Credit
		}

		// ---- Running balance ----
		balance := debit - credit
		if len(account.Splits) > 0 {
			balance += account.Splits[len(account.Splits)-1].Balance
		}

		// ---- Split ----
		account.Splits = append(account.Splits, Split{
			TransactionNumber: r.TransactionNumber,
			TransactionDate:   r.TransactionDate,
			TransactionType:   r.Type,
			TransactionMemo:   r.Memo,
			PeopleID:          r.PeopleID,
			PeopleName:        r.Name,
			Debit:             r.Debit,
			Credit:            r.Credit,
			Balance:           balance,
		})

		// ---- Totals ----
		account.TotalDebit += debit
		account.TotalCredit += credit
		account.TotalBalance = balance

		grandDebit += debit
		grandCredit += credit
	}

	// ---- Add TOTAL row per account ----
	var finalAccounts []*AccountLedger
	for _, acc := range accounts {
		finalAccounts = append(finalAccounts, acc)

		finalAccounts = append(finalAccounts, &AccountLedger{
			AccountID:     acc.AccountID,
			AccountNumber: acc.AccountNumber,
			AccountName:   "TOTAL",
			AccountType:   "",
			Splits:        []Split{},
			TotalDebit:    acc.TotalDebit,
			TotalCredit:   acc.TotalCredit,
			TotalBalance:  acc.TotalBalance,
			IsTotalRow:    true,
		})
	}

	return LedgerResponse{
		BuildingID:       buildingID,
		StartDate:        startDate,
		EndDate:          endDate,
		GrandTotalDebit:  grandDebit,
		GrandTotalCredit: grandCredit,
		Accounts:         finalAccounts,
	}
}

// GetProfitAndLossStandard
func (s *ReportService) GetProfitAndLossStandard(ctx context.Context, buildingID int, startDate string, endDate string) (*PLReport, error) {
	incomeAccounts, err := s.reportStore.GetAccountBalanceByAccountType(ctx, buildingID, startDate, endDate, "Income")
	if err != nil {
		return nil, err
	}

	expenseAccounts, err := s.reportStore.GetAccountBalanceByAccountType(ctx, buildingID, startDate, endDate, "Expense")
	if err != nil {
		return nil, err
	}

	totalExpense := 0.0
	for _, expense := range expenseAccounts {
		totalExpense += expense.Balance
	}

	totalIncome := 0.0
	for _, income := range incomeAccounts {
		totalIncome += income.Balance
	}

	// initialize with empty []
	expenseAccountsList := []store.PLAccountRow{}
	incomeAccountsList := []store.PLAccountRow{}
	for _, expense := range expenseAccounts {
		expenseAccountsList = append(expenseAccountsList, expense)
	}
	for _, income := range incomeAccounts {
		incomeAccountsList = append(incomeAccountsList, income)
	}

	return &PLReport{
		BuildingID: buildingID,
		StartDate:  startDate,
		EndDate:    endDate,
		Expenses: PLSection{
			Accounts:    expenseAccountsList,
			SectionName: "Expense",
			Total:       totalExpense,
		},
		Income: PLSection{
			Accounts:    incomeAccountsList,
			SectionName: "Income",
			Total:       totalIncome,
		},
		NetProfitLoss: totalIncome - totalExpense,
	}, nil

}

// GetProfitAndLossByUnit generates a profit and loss report grouped by unit
func (s *ReportService) GetProfitAndLossByUnit(ctx context.Context, buildingID int, startDate string, endDate string) (*dto.ProfitAndLossByUnitResponse, error) {
	// Get all units for the building
	units, err := s.unitStore.GetAll(ctx, int64(buildingID))
	if err != nil {
		return nil, err
	}

	// Get income accounts by unit
	incomeAccountsByUnit, err := s.reportStore.GetAccountBalanceByAccountTypeAndUnit(ctx, buildingID, startDate, endDate, "Income")
	if err != nil {
		return nil, err
	}

	// Get expense accounts by unit
	expenseAccountsByUnit, err := s.reportStore.GetAccountBalanceByAccountTypeAndUnit(ctx, buildingID, startDate, endDate, "Expense")
	if err != nil {
		return nil, err
	}

	// Build unit map for quick lookup
	unitMap := make(map[int]*dto.UnitColumn)
	unitMap[0] = &dto.UnitColumn{UnitID: 0, UnitName: "No Unit"} // For NULL unit_id

	for i := range units {
		unitID := int(units[i].ID)
		unitMap[unitID] = &dto.UnitColumn{
			UnitID:   unitID,
			UnitName: units[i].Name,
		}
	}

	// Get all unique unit IDs from the data
	unitIDSet := make(map[int]bool)
	for _, row := range incomeAccountsByUnit {
		if row.UnitID != nil {
			unitIDSet[*row.UnitID] = true
		} else {
			unitIDSet[0] = true
		}
	}
	for _, row := range expenseAccountsByUnit {
		if row.UnitID != nil {
			unitIDSet[*row.UnitID] = true
		} else {
			unitIDSet[0] = true
		}
	}

	// Build units list (sorted by unit name, with "No Unit" last)
	var unitsList []dto.UnitColumn
	unitColumns := make([]dto.UnitColumn, 0, len(unitIDSet))
	for id := range unitIDSet {
		if id != 0 {
			if unit, ok := unitMap[id]; ok {
				unitColumns = append(unitColumns, *unit)
			}
		}
	}
	// Sort units by name
	sort.Slice(unitColumns, func(i, j int) bool {
		return unitColumns[i].UnitName < unitColumns[j].UnitName
	})
	unitsList = append(unitsList, unitColumns...)
	// Add "No Unit" at the end if it exists
	if unitIDSet[0] {
		unitsList = append(unitsList, *unitMap[0])
	}

	// Build account maps: account_id -> map[unit_id]balance
	incomeAccountMap := make(map[int]map[int]float64)
	expenseAccountMap := make(map[int]map[int]float64)
	accountInfoMap := make(map[int]struct {
		AccountNumber int
		AccountName   string
	})

	// Process income accounts
	for _, row := range incomeAccountsByUnit {
		unitID := 0
		if row.UnitID != nil {
			unitID = *row.UnitID
		}
		if incomeAccountMap[row.AccountID] == nil {
			incomeAccountMap[row.AccountID] = make(map[int]float64)
		}
		incomeAccountMap[row.AccountID][unitID] = row.Balance
		accountInfoMap[row.AccountID] = struct {
			AccountNumber int
			AccountName   string
		}{
			AccountNumber: row.AccountNumber,
			AccountName:   row.AccountName,
		}
	}

	// Process expense accounts
	for _, row := range expenseAccountsByUnit {
		unitID := 0
		if row.UnitID != nil {
			unitID = *row.UnitID
		}
		if expenseAccountMap[row.AccountID] == nil {
			expenseAccountMap[row.AccountID] = make(map[int]float64)
		}
		expenseAccountMap[row.AccountID][unitID] = row.Balance
		accountInfoMap[row.AccountID] = struct {
			AccountNumber int
			AccountName   string
		}{
			AccountNumber: row.AccountNumber,
			AccountName:   row.AccountName,
		}
	}

	// Build income account rows
	incomeAccounts := make([]dto.AccountRow, 0)
	for accountID, balances := range incomeAccountMap {
		info := accountInfoMap[accountID]
		total := 0.0
		for _, balance := range balances {
			total += balance
		}
		incomeAccounts = append(incomeAccounts, dto.AccountRow{
			AccountID:     accountID,
			AccountNumber: info.AccountNumber,
			AccountName:   info.AccountName,
			AccountType:   "income",
			Balances:      balances,
			Total:         total,
		})
	}

	// Build expense account rows
	expenseAccounts := make([]dto.AccountRow, 0)
	for accountID, balances := range expenseAccountMap {
		info := accountInfoMap[accountID]
		total := 0.0
		for _, balance := range balances {
			total += balance
		}
		expenseAccounts = append(expenseAccounts, dto.AccountRow{
			AccountID:     accountID,
			AccountNumber: info.AccountNumber,
			AccountName:   info.AccountName,
			AccountType:   "expense",
			Balances:      balances,
			Total:         total,
		})
	}

	// Calculate totals by unit
	totalIncome := make(map[int]float64)
	totalExpenses := make(map[int]float64)
	netProfitLoss := make(map[int]float64)

	for _, account := range incomeAccounts {
		for unitID, balance := range account.Balances {
			totalIncome[unitID] += balance
		}
	}

	for _, account := range expenseAccounts {
		for unitID, balance := range account.Balances {
			totalExpenses[unitID] += balance
		}
	}

	// Calculate net profit/loss by unit
	for unitID := range totalIncome {
		netProfitLoss[unitID] = totalIncome[unitID] - totalExpenses[unitID]
	}
	for unitID := range totalExpenses {
		if _, exists := netProfitLoss[unitID]; !exists {
			netProfitLoss[unitID] = -totalExpenses[unitID]
		}
	}

	// Calculate grand totals
	grandTotalIncome := 0.0
	for _, total := range totalIncome {
		grandTotalIncome += total
	}

	grandTotalExpenses := 0.0
	for _, total := range totalExpenses {
		grandTotalExpenses += total
	}

	grandTotalNetProfitLoss := grandTotalIncome - grandTotalExpenses

	// Round values to 2 decimals
	round2 := func(v float64) float64 {
		return math.Round(v*100) / 100
	}

	for i := range incomeAccounts {
		for unitID := range incomeAccounts[i].Balances {
			incomeAccounts[i].Balances[unitID] = round2(incomeAccounts[i].Balances[unitID])
		}
		incomeAccounts[i].Total = round2(incomeAccounts[i].Total)
	}

	for i := range expenseAccounts {
		for unitID := range expenseAccounts[i].Balances {
			expenseAccounts[i].Balances[unitID] = round2(expenseAccounts[i].Balances[unitID])
		}
		expenseAccounts[i].Total = round2(expenseAccounts[i].Total)
	}

	for unitID := range totalIncome {
		totalIncome[unitID] = round2(totalIncome[unitID])
	}
	for unitID := range totalExpenses {
		totalExpenses[unitID] = round2(totalExpenses[unitID])
	}
	for unitID := range netProfitLoss {
		netProfitLoss[unitID] = round2(netProfitLoss[unitID])
	}

	return &dto.ProfitAndLossByUnitResponse{
		BuildingID:              buildingID,
		StartDate:               startDate,
		EndDate:                 endDate,
		Units:                   unitsList,
		IncomeAccounts:          incomeAccounts,
		ExpenseAccounts:         expenseAccounts,
		TotalIncome:             totalIncome,
		TotalExpenses:           totalExpenses,
		NetProfitLoss:           netProfitLoss,
		GrandTotalIncome:        round2(grandTotalIncome),
		GrandTotalExpenses:      round2(grandTotalExpenses),
		GrandTotalNetProfitLoss: round2(grandTotalNetProfitLoss),
	}, nil
}

type PLSection struct {
	Accounts    []store.PLAccountRow `json:"accounts"`
	SectionName string               `json:"section_name"`
	Total       float64              `json:"total"`
}

type PLReport struct {
	BuildingID    int       `json:"building_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	Expenses      PLSection `json:"expenses"`
	Income        PLSection `json:"income"`
	NetProfitLoss float64   `json:"net_profit_loss"`
}
