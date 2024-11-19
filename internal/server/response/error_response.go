package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var ErrInvalidJSON = errors.New("request body is invalid json")

type ValidationErrorResponse struct {
	Reason string            `json:"reason" example:"error reason"`
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"field name"`
	Message string `json:"message" example:"error message"`
}

var ValidationFailedReason = "validation failed"

func ValidationFailed(ctx echo.Context, validationErrors []ValidationError) error {
	return echo.NewHTTPError(http.StatusBadRequest, ValidationErrorResponse{
		Reason: ValidationFailedReason,
		Errors: validationErrors,
	})
}

// Format a slice of validation errors into a string. Using this function in logger.
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return "no validation errors"
	}
	var messages []string
	for _, err := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

type ErrorResponse struct {
	Reason  string `json:"reason" example:"error reason"`
	Message string `json:"message" example:"error message"`
}

var (
	BadRequestReason          = "bad request"
	UnauthorizedReason        = "unauthorized"
	ForbiddenReason           = "forbidden"
	NotFoundReason            = "not found"
	ConflictReason            = "conflict"
	InternalServerErrorReason = "internal server error"
)

func BadRequest(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
		Reason:  BadRequestReason,
		Message: err.Error(),
	})
}

func Unauthorized(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{
		Reason:  UnauthorizedReason,
		Message: err.Error(),
	})
}

func Forbidden(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusForbidden, ErrorResponse{
		Reason:  ForbiddenReason,
		Message: err.Error(),
	})
}

func NotFound(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusNotFound, ErrorResponse{
		Reason:  NotFoundReason,
		Message: err.Error(),
	})
}

func Conflict(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusConflict, ErrorResponse{
		Reason:  ConflictReason,
		Message: err.Error(),
	})
}

func InternalServerError(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, ErrorResponse{
		Reason:  InternalServerErrorReason,
		Message: err.Error(),
	})
}
