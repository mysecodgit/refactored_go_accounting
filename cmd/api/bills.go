package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

// Handlers

func (app *application) getBillsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var startDate, endDate *string
	var peopleID *int
	var status *string
	q := r.URL.Query()

	if start := q.Get("start_date"); start != "" {
		startDate = &start
	}

	if end := q.Get("end_date"); end != "" {
		endDate = &end
	}

	if pidStr := q.Get("people_id"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil {
			peopleID = &pid
		}
	}

	if statusStr := q.Get("status"); statusStr != "" {
		status = &statusStr
	}

	bills, err := app.service.Bill.GetAll(r.Context(), id, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, bills); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBillHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "billID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	bill, err := app.service.Bill.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, bill); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createBillHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBillRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.service.Bill.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Bill created successfully")
}

func (app *application) updateBillHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateBillRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	billIDStr := chi.URLParam(r, "billID")
	billID, err := strconv.ParseInt(billIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	req.ID = int(billID)

	err = app.service.Bill.Update(r.Context(), req, billID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Bill updated successfully")
}
