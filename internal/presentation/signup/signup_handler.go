package signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
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
	Name     string `json:"name" example:"Sato Taro"`
	Email    string `json:"email" example:"sato@example.com"`
	Password string `json:"password" example:"password"`
}

type SignupResponseBody struct {
	User        SignupResponseBodyUser `json:"user"`
	AccessToken string                 `json:"accessToken" example:"eyJhb..."`
}

type SignupResponseBodyUser struct {
	ID    string `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	Name  string `json:"name" example:"Sato Taro"`
	Email string `json:"email" example:"sato@example.com"`
}

// @Summary Signup
// @Description This endpoint creates a new user and issues an access token.
// @Tags Authentication API
// @Accept json
// @Produce json
// @Param body body SignupRequestBody true "Request Body"
// @Success 201 {object} SignupResponseBody
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/signup [post]
func (h *SignupHandler) Run(ctx echo.Context) error {
	req := new(SignupRequestBody)
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

	return ctx.JSON(http.StatusCreated, SignupResponseBody{
		User: SignupResponseBodyUser{
			ID:    dto.User.ID,
			Name:  dto.User.Name,
			Email: dto.User.Email,
		},
		AccessToken: dto.AccessToken,
	})
}

func (h *SignupHandler) validation(req *SignupRequestBody) (validationErrors []response.ValidationError) {
	if err := validation.ValidUserName(req.Name); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "name",
			Message: err.Error(),
		})
	}
	if err := validation.ValidUserEmail(req.Email); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "email",
			Message: err.Error(),
		})
	}
	if err := validation.ValidUserPassword(req.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "password",
			Message: err.Error(),
		})
	}
	return validationErrors
}
