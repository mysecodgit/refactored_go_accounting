package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)

/*
|---------------------------------------------------------------------------
| Interfaces
|---------------------------------------------------------------------------
*/

type JournalStore interface {
	GetAll(ctx context.Context, buildingID int64, startDate, endDate *string) ([]store.Journal, error)
	GetByID(ctx context.Context, id int64) (*store.Journal, error)
	GetByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*store.Journal, error)
	GetByTransactionID(ctx context.Context, tx *sql.Tx, transactionID int64) (*store.Journal, error)
	Create(ctx context.Context, tx *sql.Tx, j *store.Journal) (*store.Journal, error)
	Update(ctx context.Context, tx *sql.Tx, j *store.Journal) (*store.Journal, error)
	Delete(ctx context.Context, id int64) error
}

type JournalLineStore interface {
	GetAllByJournalID(ctx context.Context, journalID int64) ([]store.JournalLine, error)
	GetByID(ctx context.Context, id int64) (*store.JournalLine, error)
	Create(ctx context.Context, tx *sql.Tx, l *store.JournalLine) (*store.JournalLine, error)
	Update(ctx context.Context, tx *sql.Tx, l *store.JournalLine) (*store.JournalLine, error)
	Delete(ctx context.Context, id int64) error
	DeleteByJournalID(ctx context.Context, tx *sql.Tx, journalID int64) error
}

/*
|---------------------------------------------------------------------------
| Service
|---------------------------------------------------------------------------
*/

type JournalService struct {
	db               *sql.DB
	journalStore     JournalStore
	journalLineStore JournalLineStore
	transactionStore TransactionStore
	splitStore       SplitStore
	accountStore     AccountStore
}

/*
|---------------------------------------------------------------------------
| Constructor
|---------------------------------------------------------------------------
*/

func NewJournalService(
	db *sql.DB,
	journalStore JournalStore,
	journalLineStore JournalLineStore,
	transactionStore TransactionStore,
	splitStore SplitStore,
	accountStore AccountStore,
) *JournalService {
	return &JournalService{
		db:               db,
		journalStore:     journalStore,
		journalLineStore: journalLineStore,
		transactionStore: transactionStore,
		splitStore:       splitStore,
		accountStore:     accountStore,
	}
}

/*
|---------------------------------------------------------------------------
| Queries
|---------------------------------------------------------------------------
*/

type JournalResponseDetails struct {
	Journal      store.Journal       `json:"journal"`
	Lines []store.JournalLine `json:"lines"`
	Transaction  store.Transaction   `json:"transaction"`
	Splits  []store.Split   `json:"splits"`
}


func (s *JournalService) GetAll(ctx context.Context, buildingID int64, startDate, endDate *string) ([]store.Journal, error) {
	return s.journalStore.GetAll(ctx, buildingID, startDate, endDate)
}

func (s *JournalService) GetByID(ctx context.Context, id int64) (*JournalResponseDetails, error) {
	journal, err := s.journalStore.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("journal not found: %v", err)
	}

	lines, err := s.journalLineStore.GetAllByJournalID(ctx, journal.ID)
	if err != nil {
		return nil, fmt.Errorf("journal lines not found: %v", err)
	}

	transaction, err := s.transactionStore.GetByID(ctx, journal.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %v", err)
	}
	splits, err := s.splitStore.GetByTransactionID(ctx, transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("splits not found: %v", err)
	}

	return &JournalResponseDetails{
		Journal:      *journal,
		Lines: lines,
		Transaction:  *transaction,
		Splits: splits,
	}, nil
}

/*
|---------------------------------------------------------------------------
| Commands
|---------------------------------------------------------------------------
*/

func (s *JournalService) Create(ctx context.Context, req dto.CreateJournalRequest) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// 1. create transaction
		transaction := &store.Transaction{
			Type:              "journal",
			TransactionDate:   req.JournalDate,
			TransactionNumber: req.Reference,
			Memo:              *req.Memo,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1, // TODO: from JWT
			UnitID:            nil,
		}

		transactionID, err := s.transactionStore.Create(ctx, tx, transaction)
		if err != nil {
			fmt.Println("Error creating transaction", err)
			return err
		}

		// generate splits
		splits, err := s.GenerateJournalSplits(ctx, req.JournalPayloadDTO)
		if err != nil {
			fmt.Println("Error generating splits", err)
			return err
		}

		// validate splits
		if err := s.ValidateBalanced(splits); err != nil {
			fmt.Println("Error validating splits", err)
			return err
		}

		// create splits
		for _, split := range splits {
			split.TransactionID = *transactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				fmt.Println("Error creating splits", err)
				return err
			}
		}

		// 2. create journal
		journal := &store.Journal{
			TransactionID: *transactionID,
			Reference:     req.Reference,
			JournalDate:   req.JournalDate,
			BuildingID:    req.BuildingID,
			Memo:          req.Memo,
			TotalAmount:   &req.TotalAmount,
		}

		createdJournal, err := s.journalStore.Create(ctx, tx, journal)
		if err != nil {
			fmt.Println("Error creating journal", err)
			return err
		}

		
	

		// 3. create journal lines
		for _, line := range req.Lines {
			var debit float64 = 0.0
			var credit float64 = 0.0
			if line.Debit != nil {
				debit = *line.Debit
			}else{
				debit = 0.0
			}
			if line.Credit != nil {
				credit = *line.Credit
			}else{
				credit = 0.0
			}


			journalLine := &store.JournalLine{
				JournalID:   createdJournal.ID,
				AccountID:   int64(line.AccountID),
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Debit:       debit,
				Credit:     credit,
			}
			if _, err := s.journalLineStore.Create(ctx, tx, journalLine); err != nil {
				fmt.Println("Error creating journal lines", err)
				return err
			}
		}

		return nil
	})
}

func (s *JournalService) Update(ctx context.Context, req dto.UpdateJournalRequest, journalID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		// fetch existing journal
		existingJournal, err := s.journalStore.GetByIDTx(ctx, tx, journalID)
		if err != nil {
			fmt.Println("Error fetching existing journal", err)
			return fmt.Errorf("journal not found: %v", err)
		}

		// delete journal lines
		if err := s.journalLineStore.DeleteByJournalID(ctx, tx, journalID); err != nil {
			fmt.Println("Error deleting journal lines", err)
			return err
		}

		// update journal
		updatedJournal := &store.Journal{
			ID:            journalID,
			TransactionID: existingJournal.TransactionID,
			Reference:     req.Reference,
			JournalDate:   req.JournalDate,
			BuildingID:    req.BuildingID,
			Memo:          req.Memo,
			TotalAmount:   &req.TotalAmount,
		}

		if _, err := s.journalStore.Update(ctx, tx, updatedJournal); err != nil {
			fmt.Println("Error updating journal", err)
			return err
		}

		// recreate journal lines
		for _, line := range req.Lines {
			var debit float64 = 0.0
			var credit float64 = 0.0
			if line.Debit != nil {
				debit = *line.Debit
			}else{
				debit = 0.0
			}
			if line.Credit != nil {
				credit = *line.Credit
			}else{
				credit = 0.0
			}

			journalLine := &store.JournalLine{
				JournalID:   journalID,
				AccountID:   int64(line.AccountID),
				UnitID:      line.UnitID,
				PeopleID:    line.PeopleID,
				Description: line.Description,
				Debit:       debit,
				Credit:      credit,
			}
			if _, err := s.journalLineStore.Create(ctx, tx, journalLine); err != nil {
				fmt.Println("Error creating journal lines", err)
				return err
			}
		}

		// set old splits to 0 status

		 if err := s.splitStore.DeleteByTransactionID(ctx, tx, existingJournal.TransactionID); err != nil {
			fmt.Println("Error deleting splits", err)
			return err
		 }

		// update transaction
		transaction := &store.Transaction{
			ID:                existingJournal.TransactionID,
			Type:              "journal",
			TransactionDate:   req.JournalDate,
			TransactionNumber: req.Reference,
			Memo:              *req.Memo,
			Status:            "1",
			BuildingID:        req.BuildingID,
			UserID:            1,
			UnitID:            nil,
		}

		if _, err := s.transactionStore.Update(ctx, tx, transaction); err != nil {
			fmt.Println("Error updating transaction", err)
			return err
		}

		// re-generate splits
		splits, err := s.GenerateJournalSplits(ctx, req.JournalPayloadDTO)
		if err != nil {
			fmt.Println("Error generating splits", err)
			return err
		}

		// validate splits
		if err := s.ValidateBalanced(splits); err != nil {
			return fmt.Errorf("splits are not balanced: %v", err)
		}

		// create new splits
		for _, split := range splits {
			split.TransactionID = existingJournal.TransactionID
			if err := s.splitStore.Create(ctx, tx, &split); err != nil {
				fmt.Println("Error creating splits", err)
				return err
			}
		}

		return nil
	})
}

func (s *JournalService) GenerateJournalSplits(
	ctx context.Context,
	req dto.JournalPayloadDTO,
) ([]store.Split, error) {


	amount := req.TotalAmount
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	splits := make([]store.Split, 0)

	for _, line := range req.Lines {
		_, err := s.accountStore.GetByID(ctx, int64(line.AccountID))
		if err != nil {
			return nil, fmt.Errorf("account not found: %v", err)
		}

		splits = append(splits, store.Split{
			AccountID: int64(line.AccountID),
			Credit:    line.Credit,
			Debit:     line.Debit,
			UnitID:    line.UnitID,
			PeopleID:  line.PeopleID,
			Status:    "1",
		})

	}

	return splits, nil
}

func (s *JournalService) ValidateBalanced(splits []store.Split) error {
	var totalDebit float64 = 0.0
	var totalCredit float64 = 0.0
	for _, split := range splits {
		if split.Debit != nil {
			totalDebit += *split.Debit
		}
		if split.Credit != nil {
			totalCredit += *split.Credit
		}
	}
	if totalDebit != totalCredit {
		return fmt.Errorf("splits are not balanced")
	}
	return nil
}