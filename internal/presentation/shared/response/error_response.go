package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func BadRequest(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "BadRequest",
		Message: err.Error(),
	})
}
