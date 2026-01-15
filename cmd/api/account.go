package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

// Request structs
type createAccountRequest struct {
	AccountNumber int   `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	AccountType   int64  `json:"account_type" validate:"required"`
	BuildingID    int64  `json:"building_id" validate:"required"`
	IsDefault     int   `json:"isDefault"`
}

type updateAccountRequest struct {
	AccountNumber int   `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	AccountType   int64  `json:"account_type" validate:"required"`
	BuildingID    int64  `json:"building_id" validate:"required"`
	IsDefault     int   `json:"isDefault"`
}

// Handlers

func (app *application) getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	accounts, err := app.service.Account.GetAll(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, accounts); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	account, err := app.service.Account.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, account); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	account := &store.Account{
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		AccountType:   req.AccountType,
		BuildingID:    req.BuildingID,
		IsDefault:     req.IsDefault,
	}

	if err := app.service.Account.Create(r.Context(), account); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, account); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateAccountRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	account := &store.Account{
		ID:            id,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		AccountType:   req.AccountType,
		BuildingID:    req.BuildingID,
		IsDefault:     req.IsDefault,
	}

	if err := app.service.Account.Update(r.Context(), account); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedAccount, err := app.service.Account.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedAccount); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Account.Delete(r.Context(), id); err != nil {
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
