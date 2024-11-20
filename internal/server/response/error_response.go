package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/pocgo/pkg/strutil"
)

var ErrInvalidJSON = errors.New("request body is invalid JSON")

type ProblemDetail struct {
	Type     string `json:"type" example:"https://example.com/probs/error-title"`
	Title    string `json:"title" example:"Error title"`
	Status   int    `json:"status" example:"400"`
	Detail   string `json:"detail" example:"Error detail message"`
	Instance string `json:"instance" example:"/path/to/resource"`
}

type ValidationProblemDetail struct {
	ProblemDetail
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"field name"`
	Message string `json:"message" example:"error message"`
}

func NewProblemDetail(status int, title, detail, instance string) ProblemDetail {
	return ProblemDetail{
		Type:     fmt.Sprintf("https://example.com/probs/%s", strings.ToLower(strutil.ToKebabFromSpace(http.StatusText(status)))),
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

func ValidationFailed(ctx echo.Context, validationErrors []ValidationError) error {
	problem := ValidationProblemDetail{
		ProblemDetail: NewProblemDetail(
			http.StatusBadRequest,
			"Validation Failed",
			"one or more validation errors occurred",
			ctx.Request().URL.Path,
		),
		Errors: validationErrors,
	}
	return echo.NewHTTPError(http.StatusBadRequest, problem)
}

func BadRequest(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusBadRequest,
		"Bad Request",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusBadRequest, problem)
}

func Unauthorized(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusUnauthorized,
		"Unauthorized",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusUnauthorized, problem)
}

func Forbidden(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusForbidden,
		"Forbidden",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusForbidden, problem)
}

func NotFound(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusNotFound,
		"Not Found",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusNotFound, problem)
}

func Conflict(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusConflict,
		"Conflict",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusConflict, problem)
}

func InternalServerError(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusInternalServerError,
		"Internal Server Error",
		err.Error(),
		ctx.Request().URL.Path,
	)
	return echo.NewHTTPError(http.StatusInternalServerError, problem)
}

// Format validation errors into a single string.
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
