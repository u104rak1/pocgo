package signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
)

type SignupHandler struct {
	signupUsecase authApp.ISignupUsecase
}

func NewSignupHandler(signupUsecase authApp.ISignupUsecase) *SignupHandler {
	return &SignupHandler{
		signupUsecase: signupUsecase,
	}
}

type SignupRequestBody struct {
	User SignupRequestBodyUser `json:"user" validate:"required"`
}

type SignupRequestBodyUser struct {
	Name     string                   `json:"name" validate:"required,min=1,max=20" example:"Sato Taro"`
	Email    string                   `json:"email" validate:"required,validEmail" example:"sato@example.com"`
	Password string                   `json:"password" validate:"required,min=8,max=20" example:"password"`
	Account  SignupRequestBodyAccount `json:"account" validate:"required"`
}

type SignupRequestBodyAccount struct {
	Name     string `json:"name" validate:"required,min=1,max=10" example:"For work"`
	Password string `json:"password" validate:"required,len=4" example:"1234"`
	Currency string `json:"currency" validate:"required,oneof=JPY" example:"JPY"`
}

type SignupResponseBody struct {
	User        SignupResponseBodyUser `json:"user"`
	AccessToken string                 `json:"accessToken" example:"eyJhb..."`
}

type SignupResponseBodyUser struct {
	ID      string                    `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	Name    string                    `json:"name" example:"Sato Taro"`
	Email   string                    `json:"email" example:"sato@example.com"`
	Account SignupResponseBodyAccount `json:"account"`
}

type SignupResponseBodyAccount struct {
	ID        string  `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C7LE"`
	Name      string  `json:"name" example:"For work"`
	Balance   float64 `json:"balance" example:"0"`
	Currency  string  `json:"currency" example:"JPY"`
	UpdatedAt string  `json:"updatedAt" example:"2021-08-01T00:00:00Z"`
}

// @Summary Signup
// @Description Signup
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SignupRequestBody true "SignupRequestBody"
// @Success 201 {object} SignupResponseBody
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /v1/signup [post]
func (h *SignupHandler) Run(ctx echo.Context) error {
	req := new(SignupRequestBody)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, err)
	}

	dto, err := h.signupUsecase.Run(ctx.Request().Context(), authApp.SignupCommand{
		User: userApp.CreateUserCommand{
			Name:     req.User.Name,
			Email:    req.User.Email,
			Password: req.User.Password,
		},
		Account: accountApp.CreateAccountCommand{
			Name:     req.User.Account.Name,
			Password: req.User.Account.Password,
			Currency: req.User.Account.Currency,
		},
	})
	if err != nil {
		switch err {
		case userDomain.ErrUserEmailAlreadyExists,
			authDomain.ErrAuthenticationAlreadyExists:
			return response.Conflict(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusCreated, SignupResponseBody{
		User: SignupResponseBodyUser{
			ID:    dto.User.ID,
			Name:  dto.User.Name,
			Email: dto.User.Email,
			Account: SignupResponseBodyAccount{
				ID:        dto.Account.ID,
				Name:      dto.Account.Name,
				Balance:   dto.Account.Balance,
				Currency:  dto.Account.Currency,
				UpdatedAt: dto.Account.UpdatedAt,
			},
		},
		AccessToken: dto.AccessToken,
	})
}
