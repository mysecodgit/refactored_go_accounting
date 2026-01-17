package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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