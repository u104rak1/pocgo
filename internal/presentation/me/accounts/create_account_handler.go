package accounts

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	"github.com/ucho456job/pocgo/internal/config"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

type CreateAccountHandler struct {
	createAccountUC accountApp.ICreateAccountUsecase
}

func NewCreateAccountHandler(createAccountUsecase accountApp.ICreateAccountUsecase) *CreateAccountHandler {
	return &CreateAccountHandler{
		createAccountUC: createAccountUsecase,
	}
}

type CreateAccountRequest struct {
	Name     string `json:"name" example:"For work"`
	Password string `json:"password" example:"1234"`
	Currency string `json:"currency" example:"JPY"`
}

type CreateAccountResponse struct {
	ID        string  `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C7LE"`
	Name      string  `json:"name" example:"For work"`
	Balance   float64 `json:"balance" example:"0"`
	Currency  string  `json:"currency" example:"JPY"`
	UpdatedAt string  `json:"updatedAt" example:"2021-08-01T00:00:00Z"`
}

// @Summary Create Account
// @Description This endpoint creates a new account.
// @Tags Account API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Request Body"
// @Success 201 {object} CreateAccountResponse
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me/accounts [post]
func (h *CreateAccountHandler) Run(ctx echo.Context) error {
	req := new(CreateAccountRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	if validationErrors := h.validation(req); len(validationErrors) > 0 {
		return response.ValidationFailed(ctx, validationErrors)
	}

	userID, ok := ctx.Request().Context().Value(config.CtxUserIDKey()).(string)
	if !ok {
		return response.Unauthorized(ctx, config.ErrUserIDMissing)
	}

	dto, err := h.createAccountUC.Run(ctx.Request().Context(), accountApp.CreateAccountCommand{
		UserID:   userID,
		Name:     req.Name,
		Password: req.Password,
		Currency: req.Currency,
	})
	if err != nil {
		switch err {
		case userDomain.ErrNotFound:
			return response.NotFound(ctx, err)
		case accountDomain.ErrLimitReached:
			return response.Conflict(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusCreated, CreateAccountResponse{
		ID:        dto.ID,
		Name:      dto.Name,
		Balance:   dto.Balance,
		Currency:  dto.Currency,
		UpdatedAt: dto.UpdatedAt,
	})
}

func (h *CreateAccountHandler) validation(req *CreateAccountRequest) (validationErrors []response.ValidationError) {
	if err := validation.ValidAccountName(req.Name); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "name",
			Message: err.Error(),
		})
	}
	if err := validation.ValidAccountPassword(req.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "password",
			Message: err.Error(),
		})
	}
	if err := validation.ValidCurrency(req.Currency); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "currency",
			Message: err.Error(),
		})
	}
	return validationErrors
}
