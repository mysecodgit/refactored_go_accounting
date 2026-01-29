package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createPermissionRequest struct {
	Module string `json:"module" validate:"required"`
	Action string `json:"action" validate:"required"`
	Key    string `json:"key" validate:"required"`
}

type updatePermissionRequest struct {
	Module string `json:"module" validate:"required"`
	Action string `json:"action" validate:"required"`
	Key    string `json:"key" validate:"required"`
}

func (app *application) getPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	permissions, err := app.service.Permission.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, permissions); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPermissionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "permissionID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	permission, err := app.service.Permission.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, permission); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	var req createPermissionRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	permission := &store.Permission{
		Module: req.Module,
		Action: req.Action,
		Key:    req.Key,
	}

	if err := app.service.Permission.Create(r.Context(), permission); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, permission); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePermissionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "permissionID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updatePermissionRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	permission := &store.Permission{
		ID:     id,
		Module: req.Module,
		Action: req.Action,
		Key:    req.Key,
	}

	if err := app.service.Permission.Update(r.Context(), permission); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedPermission, err := app.service.Permission.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedPermission); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePermissionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "permissionID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Permission.Delete(r.Context(), id); err != nil {
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
