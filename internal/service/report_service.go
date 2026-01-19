package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type ReportStore interface {
	GetBalanceSheet(ctx context.Context, buildingID int, asOfDate string) ([]store.AccountBalance, error)
	GetTrialBalance(ctx context.Context, buildingID int, asOfDate string) ([]store.TrialBalanceAccount, error)
	GetCustomerBalanceSummary(ctx context.Context, buildingID int, asOfDate string) ([]store.CustomerSummary, error)
}

type ReportService struct {
	reportStore ReportStore
}

func NewReportService(
	reportStore ReportStore,
) *ReportService {
	return &ReportService{
		reportStore: reportStore,
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
			PeopleID: customer.PeopleID,
			PeopleName: customer.PeopleName,
			Balance: customer.Balance,
		})
		totalBalance += customer.Balance
	}

	return &dto.CustomerBalanceSummaryResponse{
		BuildingID: buildingID,
		AsOfDate: asOfDate,
		Customers: customersList,
		TotalBalance: totalBalance,
	}, nil
	
}