// Package orgserver ...
package orgserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	orgerror "orgstructure/internal/errors"
	"orgstructure/internal/server/middleware"
)

// Server ...
type Server struct {
	// router *http.ServeMux
	logger *slog.Logger
}

// ErrorResponse ...
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// Respond ...
func (s *Server) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	const op = "Server.Respond"

	w.WriteHeader(code)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log := s.logger.With(
				slog.String("op:", op),
				slog.String("requestID:", middleware.GetRequestIDFromRequest(r)),
			)
			log.Warn(err.Error())
		}
	}
}

// Error ...
func (s *Server) Error(w http.ResponseWriter, r *http.Request, op string, err error) {
	var resp ErrorResponse
	log := s.logger.With(
		slog.String("op", op),
		slog.String("requestID:", middleware.GetRequestIDFromRequest(r)),
	)
	log.Warn(err.Error())

	var code int
	switch {
	case errors.Is(err, orgerror.ErrEmployeeNotFound):
		code = http.StatusNotFound
		resp = ErrorResponse{
			Code:    "DEPARTMENT_NOT_FOUND",
			Message: "Department not found",
		}

	case errors.Is(err, orgerror.ErrDuplicateName):
		code = http.StatusConflict
		resp = ErrorResponse{
			Code:    "EMPLOYEE_ALREADY_EXISTS",
			Message: "Employee with this email already exists",
		}
	case errors.Is(err, orgerror.ErrInvalidDepth):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid depth",
		}
	case errors.Is(err, orgerror.ErrInvalidIncludeEmployees):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid include employees",
		}
	case errors.Is(err, orgerror.ErrInvalidID):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid ID",
		}
	case errors.Is(err, orgerror.ErrInvalidMode):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid mode",
		}
	case errors.Is(err, orgerror.ErrInvalidReassignToDepartmentID):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid reassign to department id",
		}

	default:
		code = http.StatusInternalServerError
		resp = ErrorResponse{
			Message: "Internal server error",
		}
	}

	s.Respond(w, r, code, resp)
}
