package accounts

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountApp "github.com/u104rak1/pocgo/internal/application/account"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
	"github.com/u104rak1/pocgo/internal/server/response"
)

type CreateAccountHandler struct {
	createAccountUC accountApp.ICreateAccountUsecase
}

func NewCreateAccountHandler(createAccountUsecase accountApp.ICreateAccountUsecase) *CreateAccountHandler {
	return &CreateAccountHandler{
		createAccountUC: createAccountUsecase,
	}
}

type CreateAccountRequestBody struct {
	// 3 ～ 20 文字のアカウント名
	Name string `json:"name" example:"For work"`

	// 4 桁のパスワード
	Password string `json:"password" example:"1234"`

	// 通貨（JPY または USD）
	Currency string `json:"currency" example:"JPY"`
}

type CreateAccountRequest struct {
	CreateAccountRequestBody
}

type CreateAccountResponse struct {
	// 口座ID
	ID string `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C7LE"`

	// 口座名
	Name string `json:"name" example:"For work"`

	// 口座残高
	Balance float64 `json:"balance" example:"0"`

	// 通貨
	Currency string `json:"currency" example:"JPY"`

	// 口座の更新日時
	UpdatedAt string `json:"updatedAt" example:"2021-08-01T00:00:00Z"`
}

// @Summary 口座の作成
// @Description 新しい口座を作成します。
// @Tags Account API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAccountRequestBody true "Request Body"
// @Success 201 {object} CreateAccountResponse
// @Failure 400 {object} response.ValidationProblemDetail "Validation Failed or Bad Request"
// @Failure 401 {object} response.ProblemDetail "Unauthorized"
// @Failure 404 {object} response.ProblemDetail "Not Found"
// @Failure 409 {object} response.ProblemDetail "Conflict"
// @Failure 500 {object} response.ProblemDetail "Internal Server Error"
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
			Field:   "body.name",
			Message: err.Error(),
		})
	}
	if err := validation.ValidAccountPassword(req.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "body.password",
			Message: err.Error(),
		})
	}
	if err := validation.ValidCurrency(req.Currency); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "body.currency",
			Message: err.Error(),
		})
	}
	return validationErrors
}
