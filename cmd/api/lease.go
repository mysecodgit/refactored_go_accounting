package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mysecodgit/go_accounting/internal/dto"
)

func (app *application) getLeasesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "buildingID")
	buildingID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// var startDate, endDate *string
	// var peopleID *int
	// var status *int

	// q := r.URL.Query()

	// if start := q.Get("start_date"); start != "" {
	// 	startDate = &start
	// }

	// if end := q.Get("end_date"); end != "" {
	// 	endDate = &end
	// }

	// if s := q.Get("status"); s != "" {
	// 	statusInt, err := strconv.Atoi(s)
	// 	if err != nil {
	// 		app.badRequestError(w, r, err)
	// 		return
	// 	}
	// 	status = &statusInt
	// }

	// if pidStr := q.Get("people_id"); pidStr != "" {
	// 	if pid, err := strconv.Atoi(pidStr); err == nil {
	// 		peopleID = &pid
	// 	}
	// }

	leases, err := app.service.Lease.GetAll(r.Context(), buildingID, nil, nil, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, leases)
}

func (app *application) createLeaseHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLeaseRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// validate request
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	_,err := app.service.Lease.Create(r.Context(), req, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, "lease created successfully")
}		

func (app *application) getLeaseHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "leaseID")
	leaseID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	lease, err := app.service.Lease.GetByID(r.Context(), leaseID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, lease)
}

func (app *application) updateLeaseHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateLeaseRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// validate request
	if err := Validate.Struct(req); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	_,err := app.service.Lease.Update(r.Context(), req, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "lease updated successfully")
}