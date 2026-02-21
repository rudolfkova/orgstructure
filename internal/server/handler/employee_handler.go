// Package handler ...
package handler

import (
	"encoding/json"
	"net/http"
	orgerror "orgstructure/internal/errors"
	orgserver "orgstructure/internal/server"
	"orgstructure/internal/service"
	"strconv"
	"time"
)

// EmployeeHandler ...
type EmployeeHandler struct {
	service service.EmployeeService
	server  *orgserver.Server
}

// HandleCreateEmployeeInDepartment ...
func (h *EmployeeHandler) HandleCreateEmployeeInDepartment() http.HandlerFunc {
	const op = "EmployeeHandler.handleCreateEmployeeInDepartment"

	type request struct {
		Name     string     `json:"name"`
		Position string     `json:"position"`
		HiredAt  *time.Time `json:"hired_at"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.server.Error(w, r, op, orgerror.ErrInvalidID)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		ctx := r.Context()

		emp, err := h.service.CreateEmployeeInDepartment(&ctx, req.Name, req.Position, req.HiredAt, id)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusCreated, emp)
	}
}
