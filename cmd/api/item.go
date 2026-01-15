package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/store"
)

type createItemRequest struct {
	Name           string  `json:"name" validate:"required"`
	Type           string  `json:"type" validate:"required"`
	Description    string  `json:"description" validate:"required"`

	AssetAccount   *int64  `json:"asset_account"`
	IncomeAccount  *int64  `json:"income_account"`
	COGSAccount    *int64  `json:"cogs_account"`
	ExpenseAccount *int64  `json:"expense_account"`

	OnHand     float64 `json:"on_hand" validate:"required"`
	AvgCost    float64 `json:"avg_cost" validate:"required"`
	Date       string  `json:"date" validate:"required"`
	BuildingID int64   `json:"building_id" validate:"required"`
}

type updateItemRequest = createItemRequest

func (app *application) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	items, err := app.service.Item.GetAll(r.Context(), buildingID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, items)
}


func (app *application) getItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "itemID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	item, err := app.service.Item.GetByID(r.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, item)
}

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	item := &store.Item{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		OnHand:      req.OnHand,
		AvgCost:     req.AvgCost,
		Date:        date.String(), // TODO : check this time string,
		BuildingID:  req.BuildingID,
	}

	

	if err := app.service.Item.Create(r.Context(), item); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, item)
}

func (app *application) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "itemID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var req updateItemRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	item := &store.Item{
		ID:          id,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		OnHand:      req.OnHand,
		AvgCost:     req.AvgCost,
		Date:        date.String(), // TODO : check this time string,
		BuildingID:  req.BuildingID,
	}

	if err := app.service.Item.Update(r.Context(), item); err != nil {
		if err == store.ErrNotFound {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	updated, _ := app.service.Item.GetByID(r.Context(), id)
	app.jsonResponse(w, http.StatusOK, updated)
}

func (app *application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "itemID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Item.Delete(r.Context(), id); err != nil {
		if err == store.ErrNotFound {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


