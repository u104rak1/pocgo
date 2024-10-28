package signin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

type SigninHandler struct {
	signinUC authApp.ISigninUsecase
}

func NewSigninHandler(signinUsecase authApp.ISigninUsecase) *SigninHandler {
	return &SigninHandler{
		signinUC: signinUsecase,
	}
}

type SigninRequestBody struct {
	Email    string `json:"email" example:"sato@example.com"`
	Password string `json:"password" example:"password"`
}

type SigninResponseBody struct {
	AccessToken string `json:"accessToken" example:"eyJhb..."`
}

// @Summary Signin
// @Description This endpoint authenticates the user using their email and password, and issues an access token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SigninRequestBody true "SigninRequestBody"
// @Success 201 {object} SigninResponseBody
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/signin [post]
func (h *SigninHandler) Run(ctx echo.Context) error {
	req := new(SigninRequestBody)
	if err := ctx.Bind(&req); err != nil {
		return response.BadRequest(ctx, err)
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

	return ctx.JSON(http.StatusCreated, SigninResponseBody{
		AccessToken: dto.AccessToken,
	})
}

func (h *SigninHandler) validation(req *SigninRequestBody) (validationErrors []response.ValidationError) {
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
