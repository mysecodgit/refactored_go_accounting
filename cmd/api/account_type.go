package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createAccountTypeRequest struct {
	TypeName   string `json:"typeName" validate:"required"`
	Type       string `json:"type" validate:"required"`
	SubType    string `json:"sub_type" validate:"required"`
	TypeStatus string `json:"typeStatus" validate:"required"`
}

type updateAccountTypeRequest struct {
	TypeName   string `json:"typeName" validate:"required"`
	Type       string `json:"type" validate:"required"`
	SubType    string `json:"sub_type" validate:"required"`
	TypeStatus string `json:"typeStatus" validate:"required"`
}

func (app *application) getAccountTypesHandler(w http.ResponseWriter, r *http.Request) {
	list, err := app.service.AccountType.GetAll(r.Context())
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, list); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAccountTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	at, err := app.service.AccountType.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, at); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createAccountTypeHandler(w http.ResponseWriter, r *http.Request) {
	var req createAccountTypeRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	at := &store.AccountType{
		TypeName:   req.TypeName,
		Type:       req.Type,
		SubType:    req.SubType,
		TypeStatus: req.TypeStatus,
	}

	if err := app.service.AccountType.Create(r.Context(), at); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, at); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateAccountTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateAccountTypeRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	at := &store.AccountType{
		ID:         id,
		TypeName:   req.TypeName,
		Type:       req.Type,
		SubType:    req.SubType,
		TypeStatus: req.TypeStatus,
	}

	if err := app.service.AccountType.Update(r.Context(), at); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updated, err := app.service.AccountType.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updated); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteAccountTypeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountTypeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.AccountType.Delete(r.Context(), id); err != nil {
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
