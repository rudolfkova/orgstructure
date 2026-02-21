// Package handler ...
package handler

import (
	"encoding/json"
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
			h.server.Error(w, r, op, orgerror.ErrInvalidDepartmentName)
			return
		}

		var parentIDPtr *uint
		if req.ParentID != nil {
			if *req.ParentID < 0 {
				h.server.Error(w, r, op, orgerror.ErrInvalidID)
				return
			}
			parentIDVal := uint(*req.ParentID)
			parentIDPtr = &parentIDVal
		}

		ctx := r.Context()

		dept, err := h.service.CreateDepartment(ctx, req.DepartmentName, parentIDPtr)
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
		if err != nil || id < 0 {
			h.server.Error(w, r, op, orgerror.ErrInvalidID)
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
			h.server.Error(w, r, op, orgerror.ErrInvalidDepth)
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

		tree, err := h.service.GetDepartment(ctx, depth, includeEmployees, uint(id))
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusOK, tree)
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
		if err != nil || id < 0 {
			h.server.Error(w, r, op, orgerror.ErrInvalidID)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		var parentIDPtr *uint
		if req.ParentID != nil {
			if *req.ParentID < 0 {
				h.server.Error(w, r, op, orgerror.ErrInvalidID)
				return
			}
			parentIDVal := uint(*req.ParentID)
			parentIDPtr = &parentIDVal
		}

		ctx := r.Context()

		dept, err := h.service.UpdateDepartment(ctx, uint(id), req.DepartmentName, parentIDPtr)
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
		if err != nil || id < 0 {
			h.server.Error(w, r, op, orgerror.ErrInvalidID)
			return
		}

		mode := r.URL.Query().Get("mode")

		if !isMode(mode) {
			h.server.Error(w, r, op, orgerror.ErrInvalidMode)
			return
		}

		reassignToDepartmentIDstr := r.URL.Query().Get("reassign_to_department_id")
		var reassignToDepartmentIDPtr *uint
		if reassignToDepartmentIDstr != "" {
			parsed, err := strconv.Atoi(reassignToDepartmentIDstr)
			if err != nil || parsed < 0 {
				h.server.Error(w, r, op, orgerror.ErrInvalidReassignToDepartmentID)
				return
			}
			reassignToDepartmentIDVal := uint(parsed)
			reassignToDepartmentIDPtr = &reassignToDepartmentIDVal
		}

		if mode == reassignMode && reassignToDepartmentIDPtr == nil {
			h.server.Error(w, r, op, orgerror.ErrInvalidReassignToDepartmentID)
			return
		}

		ctx := r.Context()

		err = h.service.DelDepartment(ctx, uint(id), mode, reassignToDepartmentIDPtr)
		if err != nil {
			h.server.Error(w, r, op, err)
			return
		}

		h.server.Respond(w, r, http.StatusNoContent, nil)
	}
}
