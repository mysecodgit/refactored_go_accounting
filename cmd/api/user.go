package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.service.User.GetAll(r.Context())

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, users); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
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
	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
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
		ID:       id,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
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
	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
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
