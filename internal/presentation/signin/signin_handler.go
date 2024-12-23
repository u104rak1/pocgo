package signin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
	"github.com/u104rak1/pocgo/internal/server/response"
)

type SigninHandler struct {
	signinUC authApp.ISigninUsecase
}

func NewSigninHandler(signinUsecase authApp.ISigninUsecase) *SigninHandler {
	return &SigninHandler{
		signinUC: signinUsecase,
	}
}

type SigninRequest struct {
	// ログイン用メールアドレス
	Email string `json:"email" example:"sato@example.com"`

	// ログイン用パスワード
	Password string `json:"password" example:"password"`
}

type SigninResponse struct {
	// アクセストークン
	AccessToken string `json:"accessToken" example:"eyJhb..."`
}

// @Summary サインイン
// @Description ユーザーのメールアドレスとパスワードを使用してユーザーを認証し、アクセストークンを発行します。
// @Tags Authentication API
// @Accept json
// @Produce json
// @Param request body SigninRequest true "Request Body"
// @Success 201 {object} SigninResponse
// @Failure 400 {object} response.ValidationProblemDetail "Validation Failed or Bad Request"
// @Failure 401 {object} response.ProblemDetail "Unauthorized"
// @Failure 500 {object} response.ProblemDetail "Internal Server Error"
// @Router /api/v1/signin [post]
func (h *SigninHandler) Run(ctx echo.Context) error {
	req := new(SigninRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	if validationErrors := h.validation(req); len(validationErrors) > 0 {
		return response.ValidationFailed(ctx, validationErrors)
	}

	dto, err := h.signinUC.Run(ctx.Request().Context(), authApp.SigninCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case authDomain.ErrAuthenticationFailed:
			return response.Unauthorized(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusCreated, SigninResponse{
		AccessToken: dto.AccessToken,
	})
}

func (h *SigninHandler) validation(req *SigninRequest) (validationErrors []response.ValidationError) {
	if err := validation.ValidUserEmail(req.Email); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "body.email",
			Message: err.Error(),
		})
	}
	if err := validation.ValidUserPassword(req.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "body.password",
			Message: err.Error(),
		})
	}

	return validationErrors
}
