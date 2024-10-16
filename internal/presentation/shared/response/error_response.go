package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ValidationErrorResponse struct {
	Code   string            `json:"code" example:"ErrorCode"`
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"field name"`
	Message string `json:"message" example:"error message"`
}

func ValidationFailed(ctx echo.Context, errors []ValidationError) error {
	return echo.NewHTTPError(http.StatusBadRequest, ValidationErrorResponse{
		Code:   "ValidationFailed",
		Errors: errors,
	})
}

type ErrorResponse struct {
	Code    string `json:"code" example:"ErrorCode"`
	Message string `json:"message" example:"error message"`
}

func BadRequest(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
		Code:    "BadRequest",
		Message: err.Error(),
	})
}

func Unauthorized(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{
		Code:    "Unauthorized",
		Message: err.Error(),
	})
}

func Forbidden(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusForbidden, ErrorResponse{
		Code:    "Forbidden",
		Message: err.Error(),
	})
}

func NotFound(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusNotFound, ErrorResponse{
		Code:    "NotFound",
		Message: err.Error(),
	})
}

func Conflict(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusConflict, ErrorResponse{
		Code:    "Conflict",
		Message: err.Error(),
	})
}

func InternalServerError(ctx echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, ErrorResponse{
		Code:    "InternalServerError",
		Message: err.Error(),
	})
}
