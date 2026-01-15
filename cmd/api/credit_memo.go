package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

func (app *application) getAllCreditMemoHandler(w http.ResponseWriter, r *http.Request) {
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

	creditMemos, err := app.service.CreditMemo.GetAll(r.Context(), buildingID, startDate, endDate, peopleID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, creditMemos)
}

func (app *application) getCreditMemoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "creditMemoID")
	creditMemoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	creditMemo, err := app.service.CreditMemo.GetByID(r.Context(), creditMemoID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, creditMemo)
}

func (app *application) createCreditMemoHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCreditMemoRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.CreditMemo.Create(r.Context(), req); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Credit memo created successfully")
}

func (app *application) updateCreditMemoHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCreditMemoRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.CreditMemo.Update(r.Context(), req); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	
	app.jsonResponse(w, http.StatusOK, "Credit memo updated successfully")
}
