package transactions

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

type ExecuteTransactionHandler struct {
	executeTransactionUC accountApp.IExecuteTransactionUsecase
}

func NewExecuteTransactionHandler(executeTransactionUsecase accountApp.IExecuteTransactionUsecase) *ExecuteTransactionHandler {
	return &ExecuteTransactionHandler{
		executeTransactionUC: executeTransactionUsecase,
	}
}

type ExecuteTransactionParams struct {
	AccountID string `param:"account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}

type ExecuteTransactionRequestBody struct {
	OperationType   string  `json:"operation_type" example:"deposit"`
	Amount          float64 `json:"amount" example:"1000"`
	TargetAccountID *string `json:"target_account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}

type ExecuteTransactionResponseBody struct {
	ID        string  `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	Balance   float64 `json:"balance" example:"1000"`
	Currency  string  `json:"currency" example:"JPY"`
	UpdatedAt string  `json:"updatedAt" example:"2024-03-20T15:00:00Z"`
}

// @Summary Execute Transaction
// @Description This endpoint executes a transaction (deposit, withdraw, or transfer) for the specified account.
// @Tags Account API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param request body ExecuteTransactionRequestBody true "Request Body"
// @Success 200 {object} ExecuteTransactionResponseBody
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me/accounts/{account_id}/transactions [post]
func (h *ExecuteTransactionHandler) Run(ctx echo.Context) error {
	params := new(ExecuteTransactionParams)
	if err := ctx.Bind(params); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	req := new(ExecuteTransactionRequestBody)
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

	dto, err := h.executeTransactionUC.Run(ctx.Request().Context(), accountApp.ExecuteTransactionCommand{
		UserID:          userID,
		AccountID:       params.AccountID,
		OperationType:   req.OperationType,
		Amount:          req.Amount,
		TargetAccountID: req.TargetAccountID,
	})
	if err != nil {
		switch err {
		case userDomain.ErrNotFound,
			accountDomain.ErrNotFound:
			return response.NotFound(ctx, err)
		case accountDomain.ErrInsufficientBalance,
			accountDomain.ErrInvalidOperation:
			return response.Conflict(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, ExecuteTransactionResponseBody{
		ID:        dto.ID,
		Balance:   dto.Balance,
		Currency:  dto.Currency,
		UpdatedAt: dto.UpdatedAt,
	})
}

func (h *ExecuteTransactionHandler) validation(req *ExecuteTransactionRequestBody) (validationErrors []response.ValidationError) {
	if err := validation.ValidTransactionOperationType(req.OperationType); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "operation_type",
			Message: err.Error(),
		})
	}
	if err := validation.ValidTransactionAmount(req.Amount); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "amount",
			Message: err.Error(),
		})
	}
	if req.OperationType == "transfer" {
		if req.TargetAccountID == nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "target_account_id",
				Message: "target account id is required for transfer operation",
			})
		} else if err := validation.ValidAccountID(*req.TargetAccountID); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "target_account_id",
				Message: err.Error(),
			})
		}
	}
	return validationErrors
}
