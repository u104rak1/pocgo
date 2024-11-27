package transaction

import (
	"context"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
)

type IExecuteTransactionUsecase interface {
	Run(ctx context.Context, cmd ExecuteTransactionCommand) (*ExecuteTransactionDTO, error)
}

type executeTransactionUsecase struct {
	accountRepo     accountDomain.IAccountRepository
	transactionServ transactionDomain.ITransactionService
	unitOfWork      unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction]
}

func NewExecuteTransactionUsecase(
	accountRepository accountDomain.IAccountRepository,
	transactionService transactionDomain.ITransactionService,
	unitOfWork unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction],
) IExecuteTransactionUsecase {
	return &executeTransactionUsecase{
		accountRepo:     accountRepository,
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
	account, err := u.accountRepo.FindByID(ctx, cmd.AccountID)
	if err != nil || account.UserID() != cmd.UserID {
		return nil, err
	}

	if err := account.ComparePassword(cmd.Password); err != nil {
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
			receiverAccount, err := u.accountRepo.FindByID(ctx, *cmd.ReceiverAccountID)
			if err != nil {
				if err == accountDomain.ErrNotFound {
					return nil, accountDomain.ErrReceiverNotFound
				}
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
