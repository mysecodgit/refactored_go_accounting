package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/env"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type updateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

func (app *application) getRolesHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	roles, err := app.service.Role.GetAllByOwnerID(r.Context(), ownerUserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, roles); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "roleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	role, err := app.service.Role.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Verify ownership
	if role.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, role); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	var req createRoleRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	role := &store.Role{
		OwnerUserID: ownerUserID,
		Name:        req.Name,
	}

	if err := app.service.Role.Create(r.Context(), role); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, role); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "roleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify ownership first
	existingRole, err := app.service.Role.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if existingRole.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	var req updateRoleRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	role := &store.Role{
		ID:          id,
		OwnerUserID: ownerUserID,
		Name:        req.Name,
	}

	if err := app.service.Role.Update(r.Context(), role); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedRole, err := app.service.Role.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedRole); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "roleID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Role.Delete(r.Context(), id, ownerUserID); err != nil {
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
