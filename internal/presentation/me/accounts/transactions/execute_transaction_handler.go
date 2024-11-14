package transactions

import (
	"net/http"

	"github.com/labstack/echo/v4"
	transactionApp "github.com/ucho456job/pocgo/internal/application/transaction"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	moneyVO "github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
)

type ExecuteTransactionHandler struct {
	execTransactionUC transactionApp.IExecuteTransactionUsecase
}

func NewExecuteTransactionHandler(executeTransactionUsecase transactionApp.IExecuteTransactionUsecase) *ExecuteTransactionHandler {
	return &ExecuteTransactionHandler{
		execTransactionUC: executeTransactionUsecase,
	}
}

type ExecuteTransactionRequestParams struct {
	AccountID string `param:"account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}

type ExecuteTransactionRequestBody struct {
	Password          string  `json:"password" example:"1234"`
	OperationType     string  `json:"operationType" example:"DEPOSIT"`
	Amount            float64 `json:"amount" example:"1000"`
	Currency          string  `json:"currency" example:"JPY"`
	RecieverAccountID *string `json:"recieverAccountId" example:"01J9R8AJ1Q2YDH1X9836GS9D87"`
}

type ExecuteTransactionRequest struct {
	ExecuteTransactionRequestParams
	ExecuteTransactionRequestBody
}

type ExecuteTransactionResponse struct {
	ID                string  `json:"id" example:"01J9R8AJ1Q2YDH1X9836GS9E89"`
	AccountID         string  `json:"accountId" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	RecieverAccountID *string `json:"recieverAccountId" example:"01J9R8AJ1Q2YDH1X9836GS9D87"`
	OperationType     string  `json:"operationType" example:"DEPOSIT"`
	Amount            float64 `json:"amount" example:"1000"`
	Currency          string  `json:"currency" example:"JPY"`
	TransactionAt     string  `json:"transactionAt" example:"2024-03-20T15:00:00Z"`
}

// @Summary Execute Transaction
// @Description This endpoint executes a transaction (deposit, withdraw, or transfer) for the specified account.
// @Tags Transaction API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param request body ExecuteTransactionRequestBody true "Request Body"
// @Success 200 {object} ExecuteTransactionResponse
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me/accounts/{account_id}/transactions [post]
func (h *ExecuteTransactionHandler) Run(ctx echo.Context) error {
	req := new(ExecuteTransactionRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	if validationErrors := h.validation(req); len(validationErrors) > 0 {
		return response.ValidationFailed(ctx, validationErrors)
	}

	dto, err := h.execTransactionUC.Run(ctx.Request().Context(), transactionApp.ExecuteTransactionCommand{
		AccountID:         req.AccountID,
		Password:          req.Password,
		OperationType:     req.OperationType,
		Amount:            req.Amount,
		Currency:          req.Currency,
		RecieverAccountID: req.RecieverAccountID,
	})
	if err != nil {
		switch err {
		case moneyVO.ErrDifferentCurrency:
			return response.BadRequest(ctx, err)
		case moneyVO.ErrInsufficientBalance,
			accountDomain.ErrUnmatchedPassword:
			return response.Forbidden(ctx, err)
		case accountDomain.ErrNotFound,
			accountDomain.ErrRecieverNotFound:
			return response.NotFound(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	return ctx.JSON(http.StatusCreated, ExecuteTransactionResponse{
		ID:                dto.ID,
		AccountID:         dto.AccountID,
		RecieverAccountID: dto.RecieverAccountID,
		OperationType:     dto.OperationType,
		Amount:            dto.Amount,
		Currency:          dto.Currency,
		TransactionAt:     dto.TransactionAt,
	})
}

func (h *ExecuteTransactionHandler) validation(req *ExecuteTransactionRequest) (validationErrors []response.ValidationError) {
	if err := validation.ValidULID(req.AccountID); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "account_id",
			Message: err.Error(),
		})
	}

	if err := validation.ValidAccountPassword(req.Password); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "password",
			Message: err.Error(),
		})
	}

	if err := validation.ValidTransactionOperationType(req.OperationType); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "operationType",
			Message: err.Error(),
		})
	}

	if err := validation.ValidCurrency(req.Currency); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "currency",
			Message: err.Error(),
		})
	} else {
		if err := validation.ValidAmount(req.Currency, req.Amount); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "amount",
				Message: err.Error(),
			})
		}
	}

	if req.RecieverAccountID != nil {
		if err := validation.ValidULID(*req.RecieverAccountID); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "recieverAccountId",
				Message: err.Error(),
			})
		}
	}

	return validationErrors
}
