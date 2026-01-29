package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/env"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type updateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type BuildingWithRoles struct {
	store.Building
	Roles []store.Role `json:"roles"`
}

type UserWithBuildings struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Username     string              `json:"username"`
	Phone        string              `json:"phone"`
	ParentUserID *int64              `json:"parent_user_id,omitempty"`
	Buildings    []BuildingWithRoles `json:"buildings"`
}

func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	userID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// Only fetch users that belong to the logged-in user (sub-users)
	users, err := app.service.User.GetAllByParentID(r.Context(), userID)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Fetch buildings and roles for each user
	usersWithBuildings := make([]UserWithBuildings, 0, len(users))
	for _, user := range users {
		buildings, err := app.service.UserBuilding.GetBuildingsByUserID(r.Context(), user.ID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		buildingsWithRoles := make([]BuildingWithRoles, 0, len(buildings))
		for _, building := range buildings {
			roles, err := app.service.UserBuildingRole.GetRolesByUserAndBuilding(r.Context(), user.ID, building.ID)
			if err != nil {
				app.internalServerError(w, r, err)
				return
			}

			buildingsWithRoles = append(buildingsWithRoles, BuildingWithRoles{
				Building: building,
				Roles:    roles,
			})
		}

		usersWithBuildings = append(usersWithBuildings, UserWithBuildings{
			ID:           user.ID,
			Name:         user.Name,
			Username:     user.Username,
			Phone:        user.Phone,
			ParentUserID: user.ParentUserID,
			Buildings:    buildingsWithRoles,
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, usersWithBuildings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.service.User.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Verify that the user belongs to the logged-in user
	if user.ParentUserID == nil || *user.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	// Fetch buildings and roles for the user
	buildings, err := app.service.UserBuilding.GetBuildingsByUserID(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	buildingsWithRoles := make([]BuildingWithRoles, 0, len(buildings))
	for _, building := range buildings {
		roles, err := app.service.UserBuildingRole.GetRolesByUserAndBuilding(r.Context(), user.ID, building.ID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		buildingsWithRoles = append(buildingsWithRoles, BuildingWithRoles{
			Building: building,
			Roles:    roles,
		})
	}

	userWithBuildings := UserWithBuildings{
		ID:           user.ID,
		Name:         user.Name,
		Username:     user.Username,
		Phone:        user.Phone,
		ParentUserID: user.ParentUserID,
		Buildings:    buildingsWithRoles,
	}

	if err := app.jsonResponse(w, http.StatusOK, userWithBuildings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	var req createUserRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Name:         req.Name,
		Username:     req.Username,
		Phone:        req.Phone,
		Password:     req.Password,
		ParentUserID: &parentUserID, // Set parent_user_id to logged-in user's ID
	}

	if err := app.service.User.Create(r.Context(), user); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// First verify that the user belongs to the logged-in user
	existingUser, err := app.service.User.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if existingUser.ParentUserID == nil || *existingUser.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	var req updateUserRequest

	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		ID:           id,
		Name:         req.Name,
		Username:     req.Username,
		Phone:        req.Phone,
		Password:     req.Password,
		ParentUserID: existingUser.ParentUserID, // Keep the same parent_user_id
	}

	if err := app.service.User.Update(r.Context(), user); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Fetch updated user to return
	updatedUser, err := app.service.User.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedUser); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// First verify that the user belongs to the logged-in user
	existingUser, err := app.service.User.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if existingUser.ParentUserID == nil || *existingUser.ParentUserID != parentUserID {
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized: user does not belong to you"))
		return
	}

	if err := app.service.User.Delete(r.Context(), id); err != nil {
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

type assignBuildingRequest struct {
	BuildingID int64 `json:"building_id" validate:"required"`
}

func (app *application) assignBuildingToUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify that the user belongs to the logged-in user
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

	var req assignBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify building exists
	_, err = app.service.Building.GetByID(r.Context(), req.BuildingID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.service.UserBuilding.AssignBuilding(r.Context(), userID, req.BuildingID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unassignBuildingFromUserHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// Fetch user id from jwt token
	parentUserID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idStr, 10, 64)
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

	// Verify that the user belongs to the logged-in user
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

	if err := app.service.UserBuilding.UnassignBuilding(r.Context(), userID, buildingID); err != nil {
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
