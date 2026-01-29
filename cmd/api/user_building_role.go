package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/env"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type assignRoleToUserBuildingRequest struct {
	RoleID int64 `json:"role_id" validate:"required"`
}

func (app *application) assignRoleToUserBuildingHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	buildingIDStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify user ownership
	user, err := app.service.User.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if user.ParentUserID == nil || *user.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	// Verify building exists
	_, err = app.service.Building.GetByID(r.Context(), buildingID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var req assignRoleToUserBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify role ownership
	role, err := app.service.Role.GetByID(r.Context(), req.RoleID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if role.OwnerUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: role does not belong to you"))
		return
	}

	if err := app.service.UserBuildingRole.AssignRole(r.Context(), userID, buildingID, req.RoleID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unassignRoleFromUserBuildingHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	buildingIDStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	roleIDStr := chi.URLParam(r, "roleID")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify user ownership
	user, err := app.service.User.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if user.ParentUserID == nil || *user.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	if err := app.service.UserBuildingRole.UnassignRole(r.Context(), userID, buildingID, roleID); err != nil {
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

func (app *application) getUserBuildingRolesHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	buildingIDStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(buildingIDStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify user ownership
	user, err := app.service.User.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if user.ParentUserID == nil || *user.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	roles, err := app.service.UserBuildingRole.GetRolesByUserAndBuilding(r.Context(), userID, buildingID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, roles); err != nil {
		app.internalServerError(w, r, err)
	}
}
