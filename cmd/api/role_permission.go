package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/env"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type assignPermissionRequest struct {
	PermissionID int64 `json:"permission_id" validate:"required"`
}

type setRolePermissionsRequest struct {
	PermissionIDs []int64 `json:"permission_ids" validate:"required"`
}

func (app *application) getRolePermissionsHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify role ownership
	role, err := app.service.Role.GetByID(r.Context(), roleID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if role.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	permissions, err := app.service.RolePermission.GetPermissionsByRoleID(r.Context(), roleID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, permissions); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) assignPermissionToRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify role ownership
	role, err := app.service.Role.GetByID(r.Context(), roleID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if role.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	var req assignPermissionRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.RolePermission.AssignPermission(r.Context(), roleID, req.PermissionID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unassignPermissionFromRoleHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	permissionIDStr := chi.URLParam(r, "permissionID")
	permissionID, err := strconv.ParseInt(permissionIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify role ownership
	role, err := app.service.Role.GetByID(r.Context(), roleID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if role.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	if err := app.service.RolePermission.UnassignPermission(r.Context(), roleID, permissionID); err != nil {
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

func (app *application) setRolePermissionsHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	ownerUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify role ownership
	role, err := app.service.Role.GetByID(r.Context(), roleID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if role.OwnerUserID != ownerUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	var req setRolePermissionsRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.RolePermission.SetRolePermissions(r.Context(), roleID, req.PermissionIDs); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
