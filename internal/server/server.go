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
	// ErrDuplicateEmail ...
	ErrDuplicateEmail = errors.New("employee with this email already exists")
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

	case errors.Is(err, ErrDuplicateEmail):
		code = http.StatusConflict
		resp = ErrorResponse{
			Code:    "EMPLOYEE_ALREADY_EXISTS",
			Message: "Employee with this email already exists",
		}

	default:
		code = http.StatusInternalServerError
		resp = ErrorResponse{
			Message: "Internal server error",
		}
	}

	s.Respond(w, r, code, resp)
}
