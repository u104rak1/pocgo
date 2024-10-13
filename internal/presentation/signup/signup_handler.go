package signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	signupApp "github.com/ucho456job/pocgo/internal/application/authentication/signup"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
)

type SignupHandler struct {
	signupUsecase signupApp.ISignupUsecase
}

func NewSignupHandler(signupUsecase signupApp.ISignupUsecase) *SignupHandler {
	return &SignupHandler{
		signupUsecase: signupUsecase,
	}
}

type SignupRequestBody struct {
	User SignupRequestBodyUser `json:"user" validate:"required"`
}

type SignupRequestBodyUser struct {
	Name     string                   `json:"name" validate:"required,min=1,max=20"`
	Email    string                   `json:"email" validate:"required,validEmail"`
	Password string                   `json:"password" validate:"required,min=8,max=20"`
	Account  SignupRequestBodyAccount `json:"account" validate:"required"`
}

type SignupRequestBodyAccount struct {
	Name     string `json:"name" validate:"required,min=1,max=10"`
	Password string `json:"password" validate:"required,len=4"`
	Currency string `json:"currency" validate:"required,oneof=JPY"`
}

type SignupResponseBody struct {
	User        SignupResponseBodyUser `json:"user"`
	AccessToken string                 `json:"accessToken"`
}

type SignupResponseBodyUser struct {
	ID      string                    `json:"id"`
	Name    string                    `json:"name"`
	Email   string                    `json:"email"`
	Account SignupResponseBodyAccount `json:"account"`
}

type SignupResponseBodyAccount struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	UpdatedAt string  `json:"updatedAt"`
}

func (h *SignupHandler) Run(ctx echo.Context) error {
	req := new(SignupRequestBody)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, err)
	}

	dto, err := h.signupUsecase.Run(ctx.Request().Context(), signupApp.SignupCommand{
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
		return response.BadRequest(ctx, err)
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
