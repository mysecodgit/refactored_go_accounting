package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createBuildingRequest struct {
	Name string `json:"name" validate:"required"`
}

type updateBuildingRequest struct {
	Name string `json:"name" validate:"required"`
}

func (app *application) getBuildingsHandler(w http.ResponseWriter, r *http.Request) {
	buildings, err := app.service.Building.GetAll(r.Context())
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, buildings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building, err := app.service.Building.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, building); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createBuildingHandler(w http.ResponseWriter, r *http.Request) {
	var req createBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building := &store.Building{Name: req.Name}

	if err := app.service.Building.Create(r.Context(), building); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, building); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building := &store.Building{
		ID:   id,
		Name: req.Name,
	}

	if err := app.service.Building.Update(r.Context(), building); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedBuilding, err := app.service.Building.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedBuilding); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Building.Delete(r.Context(), id); err != nil {
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
