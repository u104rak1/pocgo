package transactions

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	transactionApp "github.com/u104rak1/pocgo/internal/application/transaction"
	"github.com/u104rak1/pocgo/internal/config"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
	"github.com/u104rak1/pocgo/internal/server/response"
	"github.com/u104rak1/pocgo/pkg/timer"
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
	From           *string `query:"from" example:"20240101"`
	To             *string `query:"to" example:"20241231"`
	OperationTypes *string `query:"operation_types" example:"DEPOSIT,WITHDRAW,TRANSFER"`
	Sort           *string `query:"sort" example:"DESC"`
	Limit          *int    `query:"limit" example:"10"`
	Page           *int    `query:"page" example:"1"`
}

type ListTransactionsRequest struct {
	ListTransactionsParams
	ListTransactionsQuery
}

type ListTransactionsResponse struct {
	Total        int                           `json:"total" example:"1"`
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
// @Param account_id path string true "Account ID to be operated."
// @Param from query string false "The start date for filtering transactions (format: YYYYMMDD)."
// @Param to query string false "The end date for filtering transactions (format: YYYYMMDD)."
// @Param operationTypes query string false "Comma-separated transaction types to filter by. Valid values are DEPOSIT, WITHDRAW, and TRANSFER. If not specified, all transaction types are included."
// @Param sort query string false "The sorting order of transactions based on transactionAt. Valid values are ASC or DESC. Defaults to DESC."
// @Param limit query int false "The maximum number of transaction histories per page. Can be specified between 1 and 100."
// @Param page query int false "The page number for paginated results."
// @Success 200 {object} ListTransactionsResponse
// @Failure 400 {object} response.ValidationProblemDetail "Validation Failed or Bad Request"
// @Failure 401 {object} response.ProblemDetail "Unauthorized"
// @Failure 403 {object} response.ProblemDetail "Forbidden"
// @Failure 404 {object} response.ProblemDetail "Not Found"
// @Failure 500 {object} response.ProblemDetail "Internal Server Error"
// @Router /api/v1/me/accounts/{account_id}/transactions [get]
func (h *ListTransactionsHandler) Run(ctx echo.Context) error {
	req := new(ListTransactionsRequest)
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

	operationTypes := []string{}
	if req.OperationTypes != nil {
		operationTypes = strings.Split(*req.OperationTypes, ",")
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
		OperationTypes: operationTypes,
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

func (h *ListTransactionsHandler) validation(req *ListTransactionsRequest) (validationErrors []response.ValidationError) {
	if err := validation.ValidULID(req.AccountID); err != nil {
		validationErrors = append(validationErrors, response.ValidationError{
			Field:   "param.account_id",
			Message: err.Error(),
		})
	}
	if req.From != nil {
		if err := validation.ValidYYYYMMDD(*req.From); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.from",
				Message: err.Error(),
			})
		}
	}
	if req.To != nil {
		if err := validation.ValidYYYYMMDD(*req.To); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.to",
				Message: err.Error(),
			})
		}
	}
	if req.From != nil && req.To != nil {
		if err := validation.ValidateDateRange(*req.From, *req.To); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.from",
				Message: err.Error(),
			})
		}
	}
	if req.OperationTypes != nil {
		if err := validation.ValidTransactionOperationTypes(*req.OperationTypes); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.operation_types",
				Message: err.Error(),
			})
		}
	}
	if req.Sort != nil {
		if err := validation.ValidSort(*req.Sort); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.sort",
				Message: err.Error(),
			})
		}
	}
	if req.Limit != nil {
		if err := validation.ValidListTransactionsLimit(*req.Limit); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.limit",
				Message: err.Error(),
			})
		}
	}
	if req.Page != nil {
		if err := validation.ValidPage(*req.Page); err != nil {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   "query.page",
				Message: err.Error(),
			})
		}
	}
	return validationErrors
}
