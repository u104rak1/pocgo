package transaction

import (
	"context"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type IExecuteTransactionUsecase interface {
	Run(ctx context.Context, cmd ExecuteTransactionCommand) (*ExecuteTransactionDTO, error)
}

type executeTransactionUsecase struct {
	accountRepo     accountDomain.IAccountRepository
	transactionRepo transactionDomain.ITransactionRepository
	unitOfWork      unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction]
}

func NewExecuteTransactionUsecase(
	accountRepository accountDomain.IAccountRepository,
	transactionRepository transactionDomain.ITransactionRepository,
	unitOfWork unitofwork.IUnitOfWorkWithResult[transactionDomain.Transaction],
) IExecuteTransactionUsecase {
	return &executeTransactionUsecase{
		accountRepo:     accountRepository,
		transactionRepo: transactionRepository,
		unitOfWork:      unitOfWork,
	}
}

type ExecuteTransactionCommand struct {
	AccountID         string
	OperationType     string
	Amount            float64
	Currency          string
	RecieverAccountID *string
}

type ExecuteTransactionDTO struct {
	ID                string
	AccountID         string
	RecieverAccountID *string
	OperationType     string
	Amount            float64
	Currency          string
	TransactionAt     string
}

func (u *executeTransactionUsecase) Run(ctx context.Context, cmd ExecuteTransactionCommand) (*ExecuteTransactionDTO, error) {
	account, err := u.accountRepo.FindByID(ctx, cmd.AccountID)
	if err != nil {
		return nil, err
	}

	var transaction *transactionDomain.Transaction
	switch cmd.OperationType {
	case transactionDomain.Deposit:
		transaction, err = u.saveDepositWithTx(ctx, account, cmd)
		if err != nil {
			return nil, err
		}
	case transactionDomain.Withdraw:
		transaction, err = u.saveWithdrawWithTx(ctx, account, cmd)
		if err != nil {
			return nil, err
		}
	case transactionDomain.Transfer:
		transaction, err = u.saveTransferWithTx(ctx, account, cmd)
		if err != nil {
			return nil, err
		}
	}

	return &ExecuteTransactionDTO{
		ID:                transaction.ID(),
		AccountID:         transaction.AccountID(),
		RecieverAccountID: transaction.ReceiverAccountID(),
		OperationType:     transaction.OperationType(),
		Amount:            transaction.TransferAmount().Amount(),
		Currency:          transaction.TransferAmount().Currency(),
		TransactionAt:     transaction.TransactionAtString(),
	}, nil
}

func (u *executeTransactionUsecase) saveDepositWithTx(ctx context.Context, account *accountDomain.Account, cmd ExecuteTransactionCommand) (*transactionDomain.Transaction, error) {
	return u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
		if err := account.Deposit(cmd.Amount, cmd.Currency); err != nil {
			return nil, err
		}
		updatedAt := timer.Now()
		account.ChangeUpdatedAt(updatedAt)
		if err := u.accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}

		transactionID := ulid.New()
		transaction, err := transactionDomain.New(transactionID, account.ID(), nil, transactionDomain.Deposit,
			cmd.Amount, cmd.Currency, updatedAt)
		if err != nil {
			return nil, err
		}
		if err := u.transactionRepo.Save(ctx, transaction); err != nil {
			return nil, err
		}
		return transaction, nil
	})
}

func (u *executeTransactionUsecase) saveWithdrawWithTx(ctx context.Context, account *accountDomain.Account, cmd ExecuteTransactionCommand) (*transactionDomain.Transaction, error) {
	return u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
		if err := account.Withdraw(cmd.Amount, cmd.Currency); err != nil {
			return nil, err
		}
		updatedAt := timer.Now()
		account.ChangeUpdatedAt(updatedAt)
		if err := u.accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}

		transactionID := ulid.New()
		transaction, err := transactionDomain.New(transactionID, account.ID(), nil, transactionDomain.Withdraw,
			cmd.Amount, cmd.Currency, updatedAt)
		if err != nil {
			return nil, err
		}
		if err := u.transactionRepo.Save(ctx, transaction); err != nil {
			return nil, err
		}
		return transaction, nil
	})
}

func (u *executeTransactionUsecase) saveTransferWithTx(ctx context.Context, account *accountDomain.Account, cmd ExecuteTransactionCommand) (*transactionDomain.Transaction, error) {
	return u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*transactionDomain.Transaction, error) {
		reciverAccount, err := u.accountRepo.FindByID(ctx, *cmd.RecieverAccountID)
		if err != nil {
			if err == accountDomain.ErrNotFound {
				return nil, accountDomain.ErrRecieverNotFound
			}
			return nil, err
		}

		if err := account.Withdraw(cmd.Amount, cmd.Currency); err != nil {
			return nil, err
		}
		if err := reciverAccount.Deposit(cmd.Amount, cmd.Currency); err != nil {
			return nil, err
		}

		updatedAt := timer.Now()
		account.ChangeUpdatedAt(updatedAt)
		reciverAccount.ChangeUpdatedAt(updatedAt)

		if err := u.accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		if err := u.accountRepo.Save(ctx, reciverAccount); err != nil {
			return nil, err
		}

		transactionID := ulid.New()
		transaction, err := transactionDomain.New(transactionID, account.ID(), cmd.RecieverAccountID, transactionDomain.Transfer,
			cmd.Amount, cmd.Currency, updatedAt)
		if err != nil {
			return nil, err
		}
		if err := u.transactionRepo.Save(ctx, transaction); err != nil {
			return nil, err
		}
		return transaction, nil
	})
}
