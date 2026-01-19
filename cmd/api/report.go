package main

import (
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