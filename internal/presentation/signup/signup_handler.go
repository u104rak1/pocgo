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
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

type SignupHandler struct {
	signupUC authApp.ISignupUsecase
}

func NewSignupHandler(signupUsecase authApp.ISignupUsecase) *SignupHandler {
	return &SignupHandler{
		signupUC: signupUsecase,
	}
}

type SignupRequestBody struct {
	User SignupRequestBodyUser `json:"user"`
}

type SignupRequestBodyUser struct {
	Name     string                   `json:"name" example:"Sato Taro"`
	Email    string                   `json:"email" example:"sato@example.com"`
	Password string                   `json:"password" example:"password"`
	Account  SignupRequestBodyAccount `json:"account"`
}

type SignupRequestBodyAccount struct {
	Name     string `json:"name" example:"For work"`
	Password string `json:"password" example:"1234"`
	Currency string `json:"currency" example:"JPY"`
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
// @Description This endpoint creates a new user and account, and issues an access token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SignupRequestBody true "SignupRequestBody"
// @Success 201 {object} SignupResponseBody
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/signup [post]
func (h *SignupHandler) Run(ctx echo.Context) error {
	req := new(SignupRequestBody)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, err)
	}

	if validationErrors := h.validation(req); len(validationErrors) > 0 {
		return response.ValidationFailed(ctx, validationErrors)
	}

	dto, err := h.signupUC.Run(ctx.Request().Context(), authApp.SignupCommand{
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

func (h *SignupHandler) validation(req *SignupRequestBody) (validationErrors []response.ValidationError) {
	if err := validation.ValidUserName(req.User.Name); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.name",
			Message: err.Error(),
		})
	}
	if err := validation.ValidUserEmail(req.User.Email); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.email",
			Message: err.Error(),
		})
	}
	if err := validation.ValidUserPassword(req.User.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.password",
			Message: err.Error(),
		})
	}

	if err := validation.ValidAccountName(req.User.Account.Name); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.account.name",
			Message: err.Error(),
		})
	}
	if err := validation.ValidAccountPassword(req.User.Account.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.account.password",
			Message: err.Error(),
		})
	}
	if err := validation.ValidAccountCurrency(req.User.Account.Currency); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "user.account.currency",
			Message: err.Error(),
		})
	}
	return validationErrors
}
