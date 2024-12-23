package transaction

import (
	"context"
	"time"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IListTransactionsUsecase interface {
	Run(ctx context.Context, cmd ListTransactionsCommand) (*ListTransactionsDTO, error)
}

type listTransactionsUsecase struct {
	accountServ     accountDomain.IAccountService
	transactionServ transactionDomain.ITransactionService
}

func NewListTransactionsUsecase(
	accountService accountDomain.IAccountService,
	transactionRepository transactionDomain.ITransactionService,
) IListTransactionsUsecase {
	return &listTransactionsUsecase{
		accountServ:     accountService,
		transactionServ: transactionRepository,
	}
}

type ListTransactionsCommand struct {
	UserID         string
	AccountID      string
	From           *time.Time
	To             *time.Time
	OperationTypes []string
	Sort           *string
	Limit          *int
	Page           *int
}

type ListTransactionsDTO struct {
	Total        int
	Transactions []ListTransactionDTO
}

type ListTransactionDTO struct {
	ID                string
	AccountID         string
	ReceiverAccountID *string
	OperationType     string
	Amount            float64
	Currency          string
	TransactionAt     string
}

func (u *listTransactionsUsecase) Run(ctx context.Context, cmd ListTransactionsCommand) (*ListTransactionsDTO, error) {
	userID, err := idVO.UserIDFromString(cmd.UserID)
	if err != nil {
		return nil, err
	}

	accountID, err := idVO.AccountIDFromString(cmd.AccountID)
	if err != nil {
		return nil, err
	}

	_, err = u.accountServ.GetAndAuthorize(ctx, accountID, &userID, nil)
	if err != nil {
		return nil, err
	}

	transactions, total, err := u.transactionServ.ListWithTotal(ctx, transactionDomain.ListTransactionsParams{
		AccountID:      accountID,
		From:           cmd.From,
		To:             cmd.To,
		OperationTypes: cmd.OperationTypes,
		Sort:           cmd.Sort,
		Limit:          cmd.Limit,
		Page:           cmd.Page,
	})
	if err != nil {
		return nil, err
	}

	transactionDTOs := make([]ListTransactionDTO, len(transactions))
	for i, t := range transactions {
		if t == nil {
			continue
		}
		transactionDTOs[i] = ListTransactionDTO{
			ID:                t.IDString(),
			AccountID:         t.AccountIDString(),
			ReceiverAccountID: t.ReceiverAccountIDString(),
			OperationType:     t.OperationType(),
			Amount:            t.TransferAmount().Amount(),
			Currency:          t.TransferAmount().Currency(),
			TransactionAt:     t.TransactionAtString(),
		}
	}

	return &ListTransactionsDTO{
		Total:        total,
		Transactions: transactionDTOs,
	}, nil
}
