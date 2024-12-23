package signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
	"github.com/u104rak1/pocgo/internal/server/response"
)

type SignupHandler struct {
	signupUC authApp.ISignupUsecase
}

func NewSignupHandler(signupUsecase authApp.ISignupUsecase) *SignupHandler {
	return &SignupHandler{
		signupUC: signupUsecase,
	}
}

type SignupRequest struct {
	// 3 ～ 20 文字のユーザー名
	Name string `json:"name" example:"Sato Taro"`

	// ログイン用メールアドレス
	Email string `json:"email" example:"sato@example.com"`

	// 8 ～ 20 文字のログイン用パスワード
	Password string `json:"password" example:"password"`
}

type SignupResponse struct {
	// ユーザー
	User SignupResponseBodyUser `json:"user"`

	// アクセストークン
	AccessToken string `json:"accessToken" example:"eyJhb..."`
}

type SignupResponseBodyUser struct {
	// ユーザーID
	ID string `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`

	// ユーザー名
	Name string `json:"name" example:"Sato Taro"`

	// ユーザーのメールアドレス
	Email string `json:"email" example:"sato@example.com"`
}

// @Summary サインアップ
// @Description 新しいユーザーを作成し、アクセストークンを発行します。
// @Tags Authentication API
// @Accept json
// @Produce json
// @Param body body SignupRequest true "Request Body"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} response.ValidationProblemDetail "Validation Failed or Bad Request"
// @Failure 409 {object} response.ProblemDetail "Conflict"
// @Failure 500 {object} response.ProblemDetail "Internal Server Error"
// @Router /api/v1/signup [post]
func (h *SignupHandler) Run(ctx echo.Context) error {
	req := new(SignupRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	if validationErrors := h.validation(req); len(validationErrors) > 0 {
		return response.ValidationFailed(ctx, validationErrors)
	}

	dto, err := h.signupUC.Run(ctx.Request().Context(), authApp.SignupCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case userDomain.ErrEmailAlreadyExists,
			authDomain.ErrAlreadyExists:
			return response.Conflict(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusCreated, SignupResponse{
		User: SignupResponseBodyUser{
			ID:    dto.User.ID,
			Name:  dto.User.Name,
			Email: dto.User.Email,
		},
		AccessToken: dto.AccessToken,
	})
}

func (h *SignupHandler) validation(req *SignupRequest) (validationErrors []response.ValidationError) {
	if err := validation.ValidUserName(req.Name); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "body.name",
			Message: err.Error(),
		})
	}
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
