package transactions

import (
	"time"

	"github.com/labstack/echo/v4"
	transactionApp "github.com/ucho456job/pocgo/internal/application/transaction"
	"github.com/ucho456job/pocgo/internal/config"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/presentation/shared/response"
	"github.com/ucho456job/pocgo/pkg/timer"
)

type ListTransactionsHandler struct {
	listTransactionsUC transactionApp.IListTransactionsUsecase
}

func NewListTransactionsHandler(listTransactionsUC transactionApp.IListTransactionsUsecase) *ListTransactionsHandler {
	return &ListTransactionsHandler{
		listTransactionsUC: listTransactionsUC,
	}
}

type ListTransactionsParams struct {
	AccountID string `param:"account_id" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
}

type ListTransactionsQuery struct {
	From           *string  `query:"from" example:"2024-01-01"`
	To             *string  `query:"to" example:"2024-12-31"`
	OperationTypes []string `query:"operationTypes" example:"DEPOSIT,WITHDRAW,TRANSFER"`
	Sort           *string  `query:"sort" example:"DESC"`
	Limit          *int     `query:"limit" example:"10"`
	Page           *int     `query:"page" example:"1"`
}

type ListTransactionsRequest struct {
	ListTransactionsParams
	ListTransactionsQuery
}

type ListTransactionsResponse struct {
	Total        int                           `json:"total" example:"3"`
	Transactions []ListTransactionsTransaction `json:"transactions"`
}

type ListTransactionsTransaction struct {
	ID                string  `json:"id" example:"01J9R8AJ1Q2YDH1X9836GS9E89"`
	AccountID         string  `json:"accountId" example:"01J9R7YPV1FH1V0PPKVSB5C8FW"`
	ReceiverAccountID *string `json:"receiverAccountId" example:"01J9R8AJ1Q2YDH1X9836GS9D87"`
	OperationType     string  `json:"operationType" example:"DEPOSIT"`
	Amount            float64 `json:"amount" example:"1000"`
	Currency          string  `json:"currency" example:"JPY"`
	TransactionAt     string  `json:"transactionAt" example:"2024-03-20T15:00:00Z"`
}

// @Summary List Transactions
// @Description This endpoint retrieves the transaction history of the specified account.
// @Tags Transaction API
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param from query string false "From"
// @Param to query string false "To"
// @Param operationTypes query string false "Operation Types"
// @Param sort query string false "DESC or ASC"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} ListTransactionsResponse
// @Failure 400 {object} response.ValidationErrorResponse "Validation Failed or Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/me/accounts/{account_id}/transactions [get]
func (h *ListTransactionsHandler) Run(ctx echo.Context) error {
	req := new(ListTransactionsRequest)
	if err := ctx.Bind(req); err != nil {
		return response.BadRequest(ctx, response.ErrInvalidJSON)
	}

	// if validationErrors := h.validation(req); len(validationErrors) > 0 {
	// 	return response.ValidationFailed(ctx, validationErrors)
	// }

	userID, ok := ctx.Request().Context().Value(config.CtxUserIDKey()).(string)
	if !ok {
		return response.Unauthorized(ctx, config.ErrUserIDMissing)
	}

	var from, to *time.Time
	if req.From != nil {
		parsedFrom, err := timer.ParseYYYYMMDD(*req.From)
		if err != nil {
			return response.BadRequest(ctx, err)
		}
		from = &parsedFrom
	}
	if req.To != nil {
		parsedTo, err := timer.ParseYYYYMMDD(*req.To)
		if err != nil {
			return response.BadRequest(ctx, err)
		}
		to = &parsedTo
	}

	dto, err := h.listTransactionsUC.Run(ctx.Request().Context(), transactionApp.ListTransactionsCommand{
		UserID:         userID,
		AccountID:      req.AccountID,
		From:           from,
		To:             to,
		OperationTypes: req.OperationTypes,
		Sort:           req.Sort,
		Limit:          req.Limit,
		Page:           req.Page,
	})
	if err != nil {
		switch err {
		case accountDomain.ErrUnauthorized:
			return response.Forbidden(ctx, err)
		case accountDomain.ErrNotFound:
			return response.NotFound(ctx, err)
		default:
			return response.InternalServerError(ctx, err)
		}
	}

	transactions := make([]ListTransactionsTransaction, len(dto.Transactions))
	for i, t := range dto.Transactions {
		transactions[i] = ListTransactionsTransaction{
			ID:                t.ID,
			AccountID:         t.AccountID,
			ReceiverAccountID: t.ReceiverAccountID,
			OperationType:     t.OperationType,
			Amount:            t.Amount,
			Currency:          t.Currency,
			TransactionAt:     t.TransactionAt,
		}
	}

	return ctx.JSON(200, ListTransactionsResponse{
		Total:        dto.Total,
		Transactions: transactions,
	})
}
