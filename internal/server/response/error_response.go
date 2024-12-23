package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

func NewProblemDetail(status int, title, detail, instance, typeURL string) ProblemDetail {
	return ProblemDetail{
		Type:     typeURL,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

const (
	problemURL = "https://example.com/probs/"

	TitleValidationFailed   = "Validation Failed"
	TypeURLValidationFailed = problemURL + "validation-failed"
	DetailValidationFailed  = "one or more validation errors occurred"

	TitleBadRequest   = "Bad Request"
	TypeURLBadRequest = problemURL + "bad-request"

	TitleUnauthorized   = "Unauthorized"
	TypeURLUnauthorized = problemURL + "unauthorized"

	TitleForbidden   = "Forbidden"
	TypeURLForbidden = problemURL + "forbidden"

	TitleNotFound   = "Not Found"
	TypeURLNotFound = problemURL + "not-found"

	TitleConflict   = "Conflict"
	TypeURLConflict = problemURL + "conflict"

	TitleUnprocessableEntity   = "Unprocessable Entity"
	TypeURLUnprocessableEntity = problemURL + "unprocessable-entity"

	TitleInternalServerError   = "Internal Server Error"
	TypeURLInternalServerError = problemURL + "internal-server-error"
)

func ValidationFailed(ctx echo.Context, validationErrors []ValidationError) error {
	problem := ValidationProblemDetail{
		ProblemDetail: NewProblemDetail(
			http.StatusBadRequest,
			TitleValidationFailed,
			DetailValidationFailed,
			ctx.Request().URL.Path,
			TypeURLValidationFailed,
		),
		Errors: validationErrors,
	}
	return echo.NewHTTPError(http.StatusBadRequest, problem)
}

func BadRequest(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusBadRequest,
		TitleBadRequest,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLBadRequest,
	)
	return echo.NewHTTPError(http.StatusBadRequest, problem)
}

func Unauthorized(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusUnauthorized,
		TitleUnauthorized,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLUnauthorized,
	)
	return echo.NewHTTPError(http.StatusUnauthorized, problem)
}

func Forbidden(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusForbidden,
		TitleForbidden,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLForbidden,
	)
	return echo.NewHTTPError(http.StatusForbidden, problem)
}

func NotFound(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusNotFound,
		TitleNotFound,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLNotFound,
	)
	return echo.NewHTTPError(http.StatusNotFound, problem)
}

func Conflict(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusConflict,
		TitleConflict,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLConflict,
	)
	return echo.NewHTTPError(http.StatusConflict, problem)
}

func UnprocessableEntity(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusUnprocessableEntity,
		TitleUnprocessableEntity,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLUnprocessableEntity,
	)
	return echo.NewHTTPError(http.StatusUnprocessableEntity, problem)
}

func InternalServerError(ctx echo.Context, err error) error {
	problem := NewProblemDetail(
		http.StatusInternalServerError,
		TitleInternalServerError,
		err.Error(),
		ctx.Request().URL.Path,
		TypeURLInternalServerError,
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
