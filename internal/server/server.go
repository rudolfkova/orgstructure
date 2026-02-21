// Package orgserver ...
package orgserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"orgstructure/internal/server/middleware"
)

var (
	// ErrEmployeeNotFound ...
	ErrEmployeeNotFound = errors.New("employee not found")
	// ErrDuplicateName ...
	ErrDuplicateName = errors.New("employee with this name already exists")
	// ErrInvalidDepth ...
	ErrInvalidDepth = errors.New("invalid depth")
	// ErrInvalidIncludeEmployees ...
	ErrInvalidIncludeEmployees = errors.New("invalid include_employees")
	// ErrInvalidID ...
	ErrInvalidID = errors.New("invalid id")
	// ErrInvalidMode ...
	ErrInvalidMode = errors.New("invalid mode")
	// ErrInvalidReassignToDepartmentID ...
	ErrInvalidReassignToDepartmentID = errors.New("invalid reassign_to_department_id")
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
	case errors.Is(err, ErrEmployeeNotFound):
		code = http.StatusNotFound
		resp = ErrorResponse{
			Code:    "DEPARTMENT_NOT_FOUND",
			Message: "Department not found",
		}

	case errors.Is(err, ErrDuplicateName):
		code = http.StatusConflict
		resp = ErrorResponse{
			Code:    "EMPLOYEE_ALREADY_EXISTS",
			Message: "Employee with this email already exists",
		}
	case errors.Is(err, ErrInvalidDepth):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid depth",
		}
	case errors.Is(err, ErrInvalidIncludeEmployees):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid include employees",
		}
	case errors.Is(err, ErrInvalidID):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid ID",
		}
	case errors.Is(err, ErrInvalidMode):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT_DATA",
			Message: "Invalid mode",
		}
	case errors.Is(err, ErrInvalidReassignToDepartmentID):
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
