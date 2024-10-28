package me

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/config"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
)

type ReadMyProfileHandler struct {
	userReadUC userApp.IReadUserUsecase
}

func NewReadMyProfileHandler(userReadUsecase userApp.IReadUserUsecase) *ReadMyProfileHandler {
	return &ReadMyProfileHandler{
		userReadUC: userReadUsecase,
	}
}

type ReadMyProfileResponseBody struct {
	ID    string `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	Name  string `json:"name" example:"Sato Taro"`
	Email string `json:"email" example:"sato@example.com"`
}

// @Summary Read My Profile
// @Description This endpoint returns the profile of the authenticated user.
// @Tags User API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} ReadMyProfileResponseBody
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me [get]
func (h *ReadMyProfileHandler) Run(ctx echo.Context) error {
	userID, ok := ctx.Request().Context().Value(config.CtxUserIDKey()).(string)
	if !ok {
		return response.Unauthorized(ctx, config.ErrUserIDMissing)
	}

	dto, err := h.userReadUC.Run(ctx.Request().Context(), userApp.ReadUserCommand{
		ID: userID,
	})
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, ReadMyProfileResponseBody{
		ID:    dto.ID,
		Name:  dto.Name,
		Email: dto.Email,
	})
}
