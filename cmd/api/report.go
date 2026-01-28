package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *application) getBalanceSheetHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var asOfDate *string
	q := r.URL.Query()

	if asOfDateStr := q.Get("as_of_date"); asOfDateStr != "" {
		asOfDate = &asOfDateStr
	}

	if asOfDate == nil {
		today := time.Now().Format("2006-01-02")
		asOfDate = &today
	}

	balanceSheet, err := app.service.Report.GetBalanceSheet(r.Context(), int(id), *asOfDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, balanceSheet); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) getTrialBalanceHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var asOfDate *string
	q := r.URL.Query()
	if asOfDateStr := q.Get("as_of_date"); asOfDateStr != "" {
		asOfDate = &asOfDateStr
	}
	if asOfDate == nil {
		today := time.Now().Format("2006-01-02")
		asOfDate = &today
	}
	trialBalance, err := app.service.Report.GetTrialBalance(r.Context(), int(id), *asOfDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, trialBalance); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getCustomerBalanceSummaryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var asOfDate *string
	q := r.URL.Query()
	if asOfDateStr := q.Get("as_of_date"); asOfDateStr != "" {
		asOfDate = &asOfDateStr
	}

	if asOfDate == nil {
		today := time.Now().Format("2006-01-02")
		asOfDate = &today
	}
	customerBalanceSummary, err := app.service.Report.GetCustomerBalanceSummary(r.Context(), int(id), *asOfDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, customerBalanceSummary); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getCustomerBalanceDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var asOfDate *string
	var peopleID *int
	q := r.URL.Query()
	if asOfDateStr := q.Get("as_of_date"); asOfDateStr != "" {
		asOfDate = &asOfDateStr
	}

	if asOfDate == nil {
		today := time.Now().Format("2006-01-02")
		asOfDate = &today
	}

	if pidStr := q.Get("people_id"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil {
			peopleID = &pid
		}
	}

	customerBalanceDetail, err := app.service.Report.GetCustomerBalanceDetail(r.Context(), int(id), *asOfDate, peopleID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, customerBalanceDetail); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) getTransactionDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var startDate *string
	var endDate *string
	var unitID *int

	q := r.URL.Query()
	if startDateStr := q.Get("start_date"); startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr := q.Get("end_date"); endDateStr != "" {
		endDate = &endDateStr
	}

	if uidStr := q.Get("unit_id"); uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil {
			unitID = &uid
		}
	}

	var accountIDs []int

	if vals, ok := r.URL.Query()["account_id"]; ok {
		accountIDs = make([]int, 0, len(vals))
		for _, v := range vals {
			id, err := strconv.Atoi(v)
			if err != nil {
				app.badRequestError(w, r, err)
				return
			}
			accountIDs = append(accountIDs, id)
		}
	}

	transactionDetails, err := app.service.Report.GetTransactionDetails(r.Context(), int(id), *startDate, *endDate, accountIDs, unitID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, transactionDetails); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getProfitAndLossStandardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var startDate *string
	var endDate *string

	q := r.URL.Query()
	if startDateStr := q.Get("start_date"); startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr := q.Get("end_date"); endDateStr != "" {
		endDate = &endDateStr
	}

	profitAndLossStandard, err := app.service.Report.GetProfitAndLossStandard(r.Context(), int(id), *startDate, *endDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, profitAndLossStandard); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getProfitAndLossByUnitHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	var startDate *string
	var endDate *string

	q := r.URL.Query()
	if startDateStr := q.Get("start_date"); startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr := q.Get("end_date"); endDateStr != "" {
		endDate = &endDateStr
	}

	if startDate == nil || endDate == nil {
		app.badRequestError(w, r, fmt.Errorf("start_date and end_date are required"))
		return
	}

	profitAndLossByUnit, err := app.service.Report.GetProfitAndLossByUnit(r.Context(), int(id), *startDate, *endDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, profitAndLossByUnit); err != nil {
		app.internalServerError(w, r, err)
	}
}
