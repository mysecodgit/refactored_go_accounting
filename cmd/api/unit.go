package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createUnitRequest struct {
	Name       string `json:"name" validate:"required"`
	BuildingID int64  `json:"building_id" validate:"required"`
}

type updateUnitRequest struct {
	Name       string `json:"name" validate:"required"`
	BuildingID int64  `json:"building_id" validate:"required"`
}

func (app *application) getUnitsHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	units, err := app.service.Unit.GetAll(r.Context(), id)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, units); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUnitHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "unitID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	unit, err := app.service.Unit.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, unit); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUnitsByPeopleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	units, err := app.service.Unit.GetAllByPeopleID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, units); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createUnitHandler(w http.ResponseWriter, r *http.Request) {

	var req createUnitRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	unit := &store.Unit{
		Name:       req.Name,
		BuildingID: req.BuildingID,
	}

	if err := app.service.Unit.Create(r.Context(), unit); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, unit); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateUnitHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "unitID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateUnitRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	unit := &store.Unit{
		ID:         id,
		Name:       req.Name,
		BuildingID: req.BuildingID,
	}

	if err := app.service.Unit.Update(r.Context(), unit); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedUnit, err := app.service.Unit.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedUnit); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteUnitHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "unitID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Unit.Delete(r.Context(), id); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
