package transaction

import (
	"context"

	unitofwork "github.com/u104rak1/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
)

type IExecuteTransactionUsecase interface {
	Run(ctx context.Context, cmd ExecuteTransactionCommand) (*ExecuteTransactionDTO, error)
}

type executeTransactionUsecase struct {
	accountServ     accountDomain.IAccountService
	transactionServ transactionDomain.ITransactionService
	unitOfWork      unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction]
}

func NewExecuteTransactionUsecase(
	accountService accountDomain.IAccountService,
	transactionService transactionDomain.ITransactionService,
	unitOfWork unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction],
) IExecuteTransactionUsecase {
	return &executeTransactionUsecase{
		accountServ:     accountService,
		transactionServ: transactionService,
		unitOfWork:      unitOfWork,
	}
}

type ExecuteTransactionCommand struct {
	UserID            string
	AccountID         string
	Password          string
	OperationType     string
	Amount            float64
	Currency          string
	ReceiverAccountID *string
}

type ExecuteTransactionDTO struct {
	ID                string
	AccountID         string
	ReceiverAccountID *string
	OperationType     string
	Amount            float64
	Currency          string
	TransactionAt     string
}

func (u *executeTransactionUsecase) Run(ctx context.Context, cmd ExecuteTransactionCommand) (*ExecuteTransactionDTO, error) {
	account, err := u.accountServ.GetAndAuthorize(ctx, cmd.AccountID, &cmd.UserID, &cmd.Password)
	if err != nil {
		return nil, err
	}

	var transaction *transactionDomain.Transaction
	switch cmd.OperationType {
	case transactionDomain.Deposit:
		transaction, err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
			return u.transactionServ.Deposit(ctx, account, cmd.Amount, cmd.Currency)
		})
		if err != nil {
			return nil, err
		}
	case transactionDomain.Withdraw:
		transaction, err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
			return u.transactionServ.Withdraw(ctx, account, cmd.Amount, cmd.Currency)
		})
		if err != nil {
			return nil, err
		}
	case transactionDomain.Transfer:
		transaction, err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
			receiverAccount, err := u.accountServ.GetAndAuthorize(ctx, *cmd.ReceiverAccountID, nil, nil)
			if err != nil {
				return nil, err
			}
			return u.transactionServ.Transfer(ctx, account, receiverAccount, cmd.Amount, cmd.Currency)
		})
		if err != nil {
			return nil, err
		}
	}

	return &ExecuteTransactionDTO{
		ID:                transaction.ID(),
		AccountID:         transaction.AccountID(),
		ReceiverAccountID: transaction.ReceiverAccountID(),
		OperationType:     transaction.OperationType(),
		Amount:            transaction.TransferAmount().Amount(),
		Currency:          transaction.TransferAmount().Currency(),
		TransactionAt:     transaction.TransactionAtString(),
	}, nil
}
