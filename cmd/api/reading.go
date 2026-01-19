package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
	"github.com/mysecodgit/go_accounting/internal/store"
)



func (app *application) getReadingsByUnitHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "unitID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	readings, err := app.service.Reading.GetAllByUnitID(r.Context(), id)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, readings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getReadingsHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var status *string
	q := r.URL.Query()

	if statusStr := q.Get("status"); statusStr != "" {
		status = &statusStr
	}

	readings, err := app.service.Reading.GetAll(r.Context(), buildingID, status)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, readings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getLatestReadingHandler(w http.ResponseWriter, r *http.Request) {
	var itemID, unitID int64

	q := r.URL.Query()


	if itemIDStr := q.Get("item_id"); itemIDStr != "" {
		if iid, err := strconv.ParseInt(itemIDStr, 10, 64); err == nil {
			itemID = iid
		}
	}

	if unitIDStr := q.Get("unit_id"); unitIDStr != "" {
		if uid, err := strconv.ParseInt(unitIDStr, 10, 64); err == nil {
			unitID = uid
		}
	}

	reading, err := app.service.Reading.GetLatest(r.Context(), itemID, unitID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, reading); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createReadingHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReadingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err := app.service.Reading.Create(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getReadingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "readingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	reading, err := app.service.Reading.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, reading); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateReadingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "readingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req dto.UpdateReadingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	req.ID = int(id)
	err = app.service.Reading.Update(r.Context(), req)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}