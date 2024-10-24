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

var ValidationFailedCode = "ValidationFailed"

func ValidationFailed(ctx echo.Context, validationErrors []ValidationError) error {
	return ctx.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Code:   ValidationFailedCode,
		Errors: validationErrors,
	})
}

type ErrorResponse struct {
	Code    string `json:"code" example:"ErrorCode"`
	Message string `json:"message" example:"error message"`
}

var (
	BadRequestCode          = "BadRequest"
	UnauthorizedCode        = "Unauthorized"
	ForbiddenCode           = "Forbidden"
	NotFoundCode            = "NotFound"
	ConflictCode            = "Conflict"
	InternalServerErrorCode = "InternalServerError"
)

func BadRequest(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    BadRequestCode,
		Message: err.Error(),
	})
}

func Unauthorized(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    UnauthorizedCode,
		Message: err.Error(),
	})
}

func Forbidden(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusForbidden, ErrorResponse{
		Code:    ForbiddenCode,
		Message: err.Error(),
	})
}

func NotFound(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusNotFound, ErrorResponse{
		Code:    NotFoundCode,
		Message: err.Error(),
	})
}

func Conflict(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusConflict, ErrorResponse{
		Code:    ConflictCode,
		Message: err.Error(),
	})
}

func InternalServerError(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    InternalServerErrorCode,
		Message: err.Error(),
	})
}
