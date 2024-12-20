package transaction

import (
	"context"

	unitofwork "github.com/u104rak1/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
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
	userID, err := idVO.UserIDFromString(cmd.UserID)
	if err != nil {
		return nil, err
	}

	accountID, err := idVO.AccountIDFromString(cmd.AccountID)
	if err != nil {
		return nil, err
	}

	account, err := u.accountServ.GetAndAuthorize(ctx, accountID, &userID, &cmd.Password)
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
			receiverAccountID, err := idVO.AccountIDFromString(*cmd.ReceiverAccountID)
			if err != nil {
				return nil, err
			}
			receiverAccount, err := u.accountServ.GetAndAuthorize(ctx, receiverAccountID, nil, nil)
			if err != nil {
				return nil, err
			}
			return u.transactionServ.Transfer(ctx, account, receiverAccount, cmd.Amount, cmd.Currency)
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, transactionDomain.ErrUnsupportedType
	}

	return &ExecuteTransactionDTO{
		ID:                transaction.IDString(),
		AccountID:         transaction.AccountIDString(),
		ReceiverAccountID: transaction.ReceiverAccountIDString(),
		OperationType:     transaction.OperationType(),
		Amount:            transaction.TransferAmount().Amount(),
		Currency:          transaction.TransferAmount().Currency(),
		TransactionAt:     transaction.TransactionAtString(),
	}, nil
}
