package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

func (app *application) getSalesReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var startDate, endDate, status *string
	var peopleID *int

	q := r.URL.Query()

	if start := q.Get("start_date"); start != "" {
		startDate = &start
	}

	if end := q.Get("end_date"); end != "" {
		endDate = &end
	}

	if s := q.Get("status"); s != "" {
		status = &s
	}

	if pidStr := q.Get("people_id"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil {
			peopleID = &pid
		}
	}

	salesReceipts, err := app.service.SalesReceipt.GetAll(r.Context(), buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, salesReceipts)
}

func (app *application) createSalesReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSalesReceiptRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// validate request
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := app.service.SalesReceipt.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, "s")
}		

func (app *application) getSalesReceiptHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "salesReceiptID")
	salesReceiptID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	salesReceipt, err := app.service.SalesReceipt.GetByID(r.Context(), salesReceiptID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, salesReceipt)
}

func (app *application) updateSalesReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateSalesReceiptRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// validate request
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := app.service.SalesReceipt.Update(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "s")
}