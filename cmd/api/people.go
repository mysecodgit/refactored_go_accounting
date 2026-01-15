package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createPeopleRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	TypeID     int64  `json:"type_id" validate:"required"`
	BuildingID int64  `json:"building_id" validate:"required"`
}

type updatePeopleRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	TypeID     int64  `json:"type_id" validate:"required"`
	BuildingID int64  `json:"building_id" validate:"required"`
}

func (app *application) getPeopleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	peopleList, err := app.service.People.GetAll(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, peopleList); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPersonHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	person, err := app.service.People.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, person); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createPersonHandler(w http.ResponseWriter, r *http.Request) {
	var req createPeopleRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	person := &store.People{
		Name:       req.Name,
		Phone:      req.Phone,
		TypeID:     req.TypeID,
		BuildingID: req.BuildingID,
	}

	if err := app.service.People.Create(r.Context(), person); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, person); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePersonHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updatePeopleRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	person := &store.People{
		ID:         id,
		Name:       req.Name,
		Phone:      req.Phone,
		TypeID:     req.TypeID,
		BuildingID: req.BuildingID,
	}

	if err := app.service.People.Update(r.Context(), person); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedPerson, err := app.service.People.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedPerson); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePersonHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.People.Delete(r.Context(), id); err != nil {
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


func (app *application) getAvailableCreditsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	person, err := app.service.People.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, person); err != nil {
		app.internalServerError(w, r, err)
	}
}
