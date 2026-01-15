package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

func (app *application) createInvoicePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoicePaymentRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	payments, err := app.service.InvoicePayment.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, payments)
}


func (app *application) updateInvoicePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateInvoicePaymentRequest
	idStr := chi.URLParam(r, "invoicePaymentID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.InvoicePayment.Update(r.Context(), req, id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, nil)
	
}
func (app *application) getInvoicePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var startDate *string
	var endDate *string
	var status *string
	var peopleID *int

	q := r.URL.Query()

	if start := q.Get("start_date"); start != "" {
		startDate = &start
	}

	if end := q.Get("end_date"); end != "" {
		endDate = &end
	}


	if statusStr := q.Get("status"); statusStr != "" {
		status = &statusStr
	}

	if peopleIdStr := q.Get("people_id"); peopleIdStr != "" {
		if pid, err := strconv.Atoi(peopleIdStr); err == nil {
			peopleID = &pid
		}
	}

	

	payments, err := app.service.InvoicePayment.GetAll(r.Context(), id, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, payments)
}

func (app *application) getInvoicePaymentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "invoicePaymentID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	payment, err := app.service.InvoicePayment.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, payment)
}