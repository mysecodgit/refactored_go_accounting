package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

// Handlers

func (app *application) getJournalsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var startDate, endDate *string
	q := r.URL.Query()

	if start := q.Get("start_date"); start != "" {
		startDate = &start
	}

	if end := q.Get("end_date"); end != "" {
		endDate = &end
	}

	journals, err := app.service.Journal.GetAll(r.Context(), id, startDate, endDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, journals); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) getJournalHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "journalID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	journal, err := app.service.Journal.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, journal); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) createJournalHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJournalRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.service.Journal.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Journal created successfully")
}


func (app *application) updateJournalHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateJournalRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	journalIdStr := chi.URLParam(r, "journalID")
	journalId, err := strconv.ParseInt(journalIdStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	req.ID = int(journalId)

	err = app.service.Journal.Update(r.Context(), req, journalId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Journal updated successfully")
}