package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createPeopleTypeRequest struct {
	Title string `json:"title" validate:"required"`
}

type updatePeopleTypeRequest struct {
	Title string `json:"title" validate:"required"`
}

func (app *application) getPeopleTypesHandler(w http.ResponseWriter, r *http.Request) {
	types, err := app.service.PeopleType.GetAll(r.Context())
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, types); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPeopleTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	pt, err := app.service.PeopleType.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, pt); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createPeopleTypeHandler(w http.ResponseWriter, r *http.Request) {
	var req createPeopleTypeRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	pt := &store.PeopleType{Title: req.Title}

	if err := app.service.PeopleType.Create(r.Context(), pt); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, pt); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePeopleTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updatePeopleTypeRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	pt := &store.PeopleType{
		ID:    id,
		Title: req.Title,
	}

	if err := app.service.PeopleType.Update(r.Context(), pt); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedPT, err := app.service.PeopleType.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedPT); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePeopleTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "peopleTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.PeopleType.Delete(r.Context(), id); err != nil {
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
