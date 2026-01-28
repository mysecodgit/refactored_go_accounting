package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

// Handlers

func (app *application) getBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
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

	payments, err := app.service.BillPayment.GetAll(r.Context(), id, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, payments); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "paymentID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	payment, err := app.service.BillPayment.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, payment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBillPaymentRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	_, err := app.service.BillPayment.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Bill payment created successfully")
}

func (app *application) updateBillPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateBillPaymentRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	paymentIDStr := chi.URLParam(r, "paymentID")
	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	req.ID = paymentID

	err = app.service.BillPayment.Update(r.Context(), req, paymentID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Bill payment updated successfully")
}
