package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

// Handlers

func (app *application) getChecksHandler(w http.ResponseWriter, r *http.Request) {
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

	checks, err := app.service.Check.GetAll(r.Context(), id, startDate, endDate)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, checks); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) getCheckHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "checkID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	check, err := app.service.Check.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, check); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) createCheckHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCheckRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.service.Check.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Check created successfully")
}


func (app *application) updateCheckHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCheckRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	checkIdStr := chi.URLParam(r, "checkID")
	checkId, err := strconv.ParseInt(checkIdStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	req.ID = int(checkId)

	err = app.service.Check.Update(r.Context(), req, checkId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusOK, "Check updated successfully")
}