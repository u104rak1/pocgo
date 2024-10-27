package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ValidationErrorResponse struct {
	Reason string            `json:"reason" example:"ErrorReason"`
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"field name"`
	Message string `json:"message" example:"error message"`
}

var ValidationFailedReason = "ValidationFailed"

func ValidationFailed(ctx echo.Context, validationErrors []ValidationError) error {
	return ctx.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Reason: ValidationFailedReason,
		Errors: validationErrors,
	})
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
	return ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Reason:  BadRequestReason,
		Message: err.Error(),
	})
}

func Unauthorized(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
		Reason:  UnauthorizedReason,
		Message: err.Error(),
	})
}

func Forbidden(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusForbidden, ErrorResponse{
		Reason:  ForbiddenReason,
		Message: err.Error(),
	})
}

func NotFound(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusNotFound, ErrorResponse{
		Reason:  NotFoundReason,
		Message: err.Error(),
	})
}

func Conflict(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusConflict, ErrorResponse{
		Reason:  ConflictReason,
		Message: err.Error(),
	})
}

func InternalServerError(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Reason:  InternalServerErrorReason,
		Message: err.Error(),
	})
}
