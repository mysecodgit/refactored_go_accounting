package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

func (app *application) getInvoicesHandler(w http.ResponseWriter, r *http.Request) {
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

	invoices, err := app.service.Invoice.GetAll(r.Context(), buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, invoices)
}

func (app *application) getInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	_, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	invoiceIdStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(invoiceIdStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	invoice, err := app.service.Invoice.GetByID(r.Context(), invoiceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, invoice)
}

func (app *application) getPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	payments, err := app.service.Invoice.GetPayments(r.Context(), invoiceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, payments)
}
func (app *application) previewInvoiceSplitsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.InvoicePayloadDTO
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	splits, err := app.service.Invoice.PreviewInvoiceSplits(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, splits)
}

func (app *application) createInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoiceRequestDTO
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Invoice.Create(r.Context(), req); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Invoice created successfully")
}

func (app *application) updateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateInvoiceRequestDTO
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Invoice.Update(r.Context(), req); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Invoice updated successfully")
}

func (app *application) getInvoiceDiscountsHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	discounts, err := app.service.Invoice.GetInvoiceDiscounts(r.Context(), invoiceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, discounts)
}

func (app *application) applyInvoiceDiscountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoiceAppliedDiscountRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	discount, err := app.service.Invoice.CreateInvoiceDiscount(r.Context(), invoiceId, req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, discount)
}

func (app *application) getInvoiceAppliedCreditsHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	credits, err := app.service.Invoice.GetAppliedCredits(r.Context(), invoiceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, credits)
}

func (app *application) getInvoiceAvailableCreditsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	availableCredits, err := app.service.Invoice.GetInvoiceAvailableCredits(r.Context(), invoiceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, availableCredits)
}

func (app *application) applyInvoiceCreditHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoiceAppliedCreditRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "invoiceID")
	invoiceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	req.InvoiceID = int(invoiceId)

	err = app.service.Invoice.ApplyInvoiceCredits(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Invoice credit applied successfully")
}