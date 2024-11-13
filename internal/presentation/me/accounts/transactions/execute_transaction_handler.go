package transactions

import (
	"net/http"

	"github.com/labstack/echo/v4"
	transactionApp "github.com/ucho456job/pocgo/internal/application/transaction"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	moneyDomain "github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
)

type ExecuteTransactionHandler struct {
	execTransactionUC transactionApp.IExecuteTransactionUsecase
}

func NewExecuteTransactionHandler(executeTransactionUsecase transactionApp.IExecuteTransactionUsecase) *ExecuteTransactionHandler {
	return &ExecuteTransactionHandler{
		execTransactionUC: executeTransactionUsecase,
	}
}

type ExecuteTransactionRequest struct {
	AccountID         string  `param:"account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	OperationType     string  `json:"operationType" example:"deposit"`
	Amount            float64 `json:"amount" example:"1000"`
	Currency          string  `json:"currency" example:"JPY"`
	RecieverAccountID *string `json:"recieverAccountId" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}

type ExecuteTransactionResponse struct {
	ID                string  `json:"id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	AccountID         string  `json:"accountId" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	RecieverAccountID *string `json:"recieverAccountId" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	OperationType     string  `json:"operationType" example:"deposit"`
	Amount            float64 `json:"amount" example:"1000"`
	Currency          string  `json:"currency" example:"JPY"`
	TransactionAt     string  `json:"transactionAt" example:"2024-03-20T15:00:00Z"`
}

// @Summary Execute Transaction
// @Description This endpoint executes a transaction (deposit, withdraw, or transfer) for the specified account.
// @Tags Account API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param request body ExecuteTransactionRequest true "Request Body"
// @Success 200 {object} ExecuteTransactionResponse
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 409 {object} response.ErrorResponse "Conflict"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me/accounts/{account_id}/transactions [post]
func (h *ExecuteTransactionHandler) Run(ctx echo.Context) error {
	req := new(ExecuteTransactionRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	// if validationErrors := h.validation(req); len(validationErrors) > 0 {
	// 	return response.ValidationFailed(ctx, validationErrors)
	// }

	dto, err := h.execTransactionUC.Run(ctx.Request().Context(), transactionApp.ExecuteTransactionCommand{
		AccountID:         req.AccountID,
		OperationType:     req.OperationType,
		Amount:            req.Amount,
		Currency:          req.Currency,
		RecieverAccountID: req.RecieverAccountID,
	})
	if err != nil {
		switch err {
		case moneyDomain.ErrDifferentCurrency:
			return response.BadRequest(ctx, err)
		case moneyDomain.ErrInsufficientBalance:
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

// func (h *ExecuteTransactionHandler) validation(req *ExecuteTransactionRequestBody) (validationErrors []response.ValidationError) {
// 	if err := validation.ValidTransactionOperationType(req.OperationType); err != nil {
// 		validationErrors = append(validationErrors, response.ValidationError{
// 			Field:   "operation_type",
// 			Message: err.Error(),
// 		})
// 	}
// 	if err := validation.ValidTransactionAmount(req.Amount); err != nil {
// 		validationErrors = append(validationErrors, response.ValidationError{
// 			Field:   "amount",
// 			Message: err.Error(),
// 		})
// 	}
// 	if req.OperationType == "transfer" {
// 		if req.RecieverAccountID == nil {
// 			validationErrors = append(validationErrors, response.ValidationError{
// 				Field:   "target_account_id",
// 				Message: "target account id is required for transfer operation",
// 			})
// 		} else if err := validation.ValidAccountID(*req.RecieverAccountID); err != nil {
// 			validationErrors = append(validationErrors, response.ValidationError{
// 				Field:   "target_account_id",
// 				Message: err.Error(),
// 			})
// 		}
// 	}
// 	return validationErrors
// }
