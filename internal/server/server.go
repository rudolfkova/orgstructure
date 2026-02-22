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
	Logger *slog.Logger
}

// NewServer ...
func NewServer(logger *slog.Logger) *Server {
	return &Server{Logger: logger}
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
			log := s.Logger.With(
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
	log := s.Logger.With(
		slog.String("op", op),
		slog.String("requestID:", middleware.GetRequestIDFromRequest(r)),
	)
	log.Warn(err.Error())

	var code int
	switch {
	case errors.Is(err, orgerror.ErrDepartmentNotFound),
		errors.Is(err, orgerror.ErrEmployeeNotFound),
		errors.Is(err, orgerror.ErrParentNotFound):
		code = http.StatusNotFound
		resp = ErrorResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		}

	case errors.Is(err, orgerror.ErrDuplicateName):
		code = http.StatusConflict
		resp = ErrorResponse{
			Code:    "DUPLICATE_NAME",
			Message: err.Error(),
		}

	case errors.Is(err, orgerror.ErrCyclicDependency):
		code = http.StatusConflict
		resp = ErrorResponse{
			Code:    "CYCLIC_DEPENDENCY",
			Message: err.Error(),
		}

	case errors.Is(err, orgerror.ErrInvalidDepartmentName),
		errors.Is(err, orgerror.ErrInvalidFullName),
		errors.Is(err, orgerror.ErrInvalidPosition),
		errors.Is(err, orgerror.ErrInvalidDepth),
		errors.Is(err, orgerror.ErrInvalidIncludeEmployees),
		errors.Is(err, orgerror.ErrInvalidID),
		errors.Is(err, orgerror.ErrInvalidMode),
		errors.Is(err, orgerror.ErrReassignTargetInSubtree),
		errors.Is(err, orgerror.ErrInvalidReassignToDepartmentID):
		code = http.StatusBadRequest
		resp = ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		}

	default:
		code = http.StatusInternalServerError
		resp = ErrorResponse{
			Message: "Internal server error",
		}
	}

	s.Respond(w, r, code, resp)
}
