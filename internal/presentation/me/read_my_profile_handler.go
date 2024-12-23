package me

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userApp "github.com/u104rak1/pocgo/internal/application/user"
	"github.com/u104rak1/pocgo/internal/config"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/server/response"
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
	// ユーザーのID
	ID string `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`

	// ユーザーの名前
	Name string `json:"name" example:"Sato Taro"`

	// ユーザーのメールアドレス
	Email string `json:"email" example:"sato@example.com"`
}

// @Summary プロフィールの取得
// @Description 認証済みのユーザーのプロフィールを返します。
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
