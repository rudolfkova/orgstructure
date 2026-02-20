// Package handler ...
package handler

import (
	"encoding/json"
	"net/http"
	"orgstructure/internal/repository"
	orgserver "orgstructure/internal/server"
	"time"
)

// EmployeeHandler ...
type EmployeeHandler struct {
	empRepo repository.EmployeeRepository
}

// HandleCreateEmployeeInDepartment ...
func (h *EmployeeHandler) HandleCreateEmployeeInDepartment(s *orgserver.Server) http.HandlerFunc {
	const op = "EmployeeHandler.handleCreateEmployeeInDepartment"

	type request struct {
		Name     string    `json:"name"`
		Position string    `json:"position"`
		HiredAt  time.Time `json:"hired_at"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.Error(w, r, op, err)
			return
		}

		ctx := r.Context()

		emp, err := h.empRepo.CreateEmployeeInDepartment(&ctx, req.Name, req.Position, req.HiredAt)
		if err != nil {
			s.Error(w, r, op, err)
			return
		}

		s.Respond(w, r, http.StatusCreated, emp)
	}
}
