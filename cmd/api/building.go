package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mysecodgit/go_accounting/internal/env"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createBuildingRequest struct {
	Name string `json:"name" validate:"required"`
}

type updateBuildingRequest struct {
	Name string `json:"name" validate:"required"`
}

func getUserIDFromJWT(r *http.Request, jwtSecret string) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("no token provided")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("sub claim missing")
	}

	// MapClaims stores numbers as float64
	userIDFloat, ok := sub.(float64)
	if !ok {
		return 0, errors.New("invalid sub claim type")
	}

	return int64(userIDFloat), nil
}

func (app *application) getBuildingsHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me")
	// fetch user id from jwt token
	userID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	fmt.Println("userID", userID)
	fmt.Println("jwtSecret", jwtSecret)

	buildings, err := app.service.Building.GetAllByUserID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, buildings); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building, err := app.service.Building.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, building); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createBuildingHandler(w http.ResponseWriter, r *http.Request) {
	jwtSecret := env.GetString("JWT_SECRET", "dev_secret_change_me") // TODO : this is only for development, so we need to remove it in production
	// fetch user id from jwt token
	userID, err := getUserIDFromJWT(r, jwtSecret)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var req createBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building := &store.Building{Name: req.Name}

	if err := app.service.Building.Create(r.Context(), building, userID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, building); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateBuildingRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	building := &store.Building{
		ID:   id,
		Name: req.Name,
	}

	if err := app.service.Building.Update(r.Context(), building); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	updatedBuilding, err := app.service.Building.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedBuilding); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteBuildingHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Building.Delete(r.Context(), id); err != nil {
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
