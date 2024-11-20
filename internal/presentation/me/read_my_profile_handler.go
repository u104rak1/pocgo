package me

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/config"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/server/response"
)

type ReadMyProfileHandler struct {
	readUserUC userApp.IReadUserUsecase
}

func NewReadMyProfileHandler(userReadUsecase userApp.IReadUserUsecase) *ReadMyProfileHandler {
	return &ReadMyProfileHandler{
		readUserUC: userReadUsecase,
	}
}

type ReadMyProfileResponse struct {
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
// @Success 200 {object} ReadMyProfileResponse
// @Failure 401 {object} response.ProblemDetail "Unauthorized"
// @Failure 404 {object} response.ProblemDetail "Not Found"
// @Failure 500 {object} response.ProblemDetail "Internal Server Error"
// @Router /api/v1/me [get]
func (h *ReadMyProfileHandler) Run(ctx echo.Context) error {
	userID, ok := ctx.Request().Context().Value(config.CtxUserIDKey()).(string)
	if !ok {
		return response.Unauthorized(ctx, config.ErrUserIDMissing)
	}

	dto, err := h.readUserUC.Run(ctx.Request().Context(), userApp.ReadUserCommand{
		ID: userID,
	})
	if err != nil {
		switch err {
		case userDomain.ErrNotFound:
			return response.NotFound(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, ReadMyProfileResponse{
		ID:    dto.ID,
		Name:  dto.Name,
		Email: dto.Email,
	})
}
