// Package handler ...
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	orgerror "orgstructure/internal/errors"
	orgserver "orgstructure/internal/server"
	"orgstructure/internal/service"
	"strconv"
	"strings"
)

// DepartmentHandler ...
type DepartmentHandler struct {
	service service.DepartmentService
	server  *orgserver.Server
}

var (
	cascadeMode  = "cascade"
	reassignMode = "reassign"
)

func isMode(str string) bool {
	return str == cascadeMode || str == reassignMode
}

// HandleCreateDepartment ...
func (h *DepartmentHandler) HandleCreateDepartment() http.HandlerFunc {
	const op = "EmployeeHandler.HandleCreateDepartment"

	type request struct {
		DepartmentName string `json:"name"`
		ParentID       *int   `json:"parent_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		if strings.TrimSpace(req.DepartmentName) == "" {
			h.server.Error(w, r, op, fmt.Errorf("invalid department_name"))
			return
		}

		ctx := r.Context()

		dept, err := h.service.CreateDepartment(ctx, req.DepartmentName, req.ParentID)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusCreated, dept)
	}
}

// HandleGetDepartment ...
func (h *DepartmentHandler) HandleGetDepartment() http.HandlerFunc {
	const op = "EmployeeHandler.HandleGetDepartment"

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.server.Error(w, r, op, fmt.Errorf("invalid id"))
			return
		}
		depthStr := r.URL.Query().Get("depth")
		includeStr := r.URL.Query().Get("include_employees")

		depth := 1 // default
		if depthStr != "" {
			parsed, err := strconv.Atoi(depthStr)
			if err != nil {
				h.server.Error(w, r, op, orgerror.ErrInvalidDepth)
				return
			}
			depth = parsed
		}
		if depth < 0 {
			h.server.Error(w, r, op, fmt.Errorf("invalid depth"))
		}

		includeEmployees := true // default
		if includeStr != "" {
			parsed, err := strconv.ParseBool(includeStr)
			if err != nil {
				h.server.Error(w, r, op, orgerror.ErrInvalidIncludeEmployees)
				return
			}
			includeEmployees = parsed
		}

		ctx := r.Context()

		dept, err := h.service.GetDepartment(ctx, depth, includeEmployees, id)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusOK, dept)
	}
}

// HandleUpdateDepartment ...
func (h *DepartmentHandler) HandleUpdateDepartment() http.HandlerFunc {
	const op = "EmployeeHandler.HandleUpdateDepartment"

	type request struct {
		DepartmentName *string `json:"name"`
		ParentID       *int    `json:"parent_id"`
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

		dept, err := h.service.UpdateDepartment(ctx, req.DepartmentName, req.ParentID, id)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusOK, dept)
	}
}

// HandleDelDepartment ...
func (h *DepartmentHandler) HandleDelDepartment() http.HandlerFunc {
	const op = "EmployeeHandler.HandleDelDepartment"

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.server.Error(w, r, op, orgerror.ErrInvalidID)
			return
		}

		mode := r.URL.Query().Get("mode")

		if !isMode(mode) {
			h.server.Error(w, r, op, orgerror.ErrInvalidMode)
			return
		}

		reassignToDepartmentIDstr := r.URL.Query().Get("reassign_to_department_id")
		var reassignToDepartmentID *int
		if reassignToDepartmentIDstr != "" {
			parsed, err := strconv.Atoi(reassignToDepartmentIDstr)
			if err != nil {
				h.server.Error(w, r, op, orgerror.ErrInvalidReassignToDepartmentID)
				return
			}
			reassignToDepartmentID = &parsed
		}

		if mode == reassignMode && reassignToDepartmentID == nil {
			h.server.Error(w, r, op, orgerror.ErrInvalidReassignToDepartmentID)
			return
		}

		ctx := r.Context()

		err = h.service.DelDepartment(ctx, mode, reassignToDepartmentID, id)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusNoContent, nil)
	}
}
