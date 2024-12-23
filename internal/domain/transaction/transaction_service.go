package transaction

import (
	"context"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/pkg/timer"
)

type ITransactionService interface {
	Deposit(ctx context.Context, account *accountDomain.Account, amount float64, currency string) (*Transaction, error)
	Withdrawal(ctx context.Context, account *accountDomain.Account, amount float64, currency string) (*Transaction, error)
	Transfer(ctx context.Context, senderAccount *accountDomain.Account, receiverAccount *accountDomain.Account, amount float64, currency string) (*Transaction, error)
	ListWithTotal(ctx context.Context, params ListTransactionsParams) (transactions []*Transaction, total int, err error)
}

type transactionService struct {
	accountRepo     accountDomain.IAccountRepository
	transactionRepo ITransactionRepository
}

func NewService(
	accountRepository accountDomain.IAccountRepository,
	transactionRepository ITransactionRepository) ITransactionService {
	return &transactionService{
		accountRepo:     accountRepository,
		transactionRepo: transactionRepository,
	}
}

func (s *transactionService) Deposit(
	ctx context.Context,
	account *accountDomain.Account,
	amount float64,
	currency string,
) (*Transaction, error) {
	if err := account.Deposit(amount, currency); err != nil {
		return nil, err
	}
	updatedAt := timer.Now()
	account.ChangeUpdatedAt(updatedAt)
	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	transaction, err := New(account.ID(), nil, Deposit, amount, currency, updatedAt)
	if err != nil {
		return nil, err
	}
	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) Withdrawal(
	ctx context.Context,
	account *accountDomain.Account,
	amount float64,
	currency string,
) (*Transaction, error) {
	if err := account.Withdrawal(amount, currency); err != nil {
		return nil, err
	}
	updatedAt := timer.Now()
	account.ChangeUpdatedAt(updatedAt)
	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	transaction, err := New(account.ID(), nil, Withdrawal, amount, currency, updatedAt)
	if err != nil {
		return nil, err
	}
	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) Transfer(
	ctx context.Context,
	senderAccount *accountDomain.Account,
	receiverAccount *accountDomain.Account,
	amount float64,
	currency string,
) (*Transaction, error) {
	if err := receiverAccount.Deposit(amount, currency); err != nil {
		return nil, err
	}
	if err := senderAccount.Withdrawal(amount, currency); err != nil {
		return nil, err
	}

	updatedAt := timer.Now()
	senderAccount.ChangeUpdatedAt(updatedAt)
	receiverAccount.ChangeUpdatedAt(updatedAt)

	if err := s.accountRepo.Save(ctx, senderAccount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Save(ctx, receiverAccount); err != nil {
		return nil, err
	}

	receiverAccountID := receiverAccount.ID()
	transaction, err := New(senderAccount.ID(), &receiverAccountID, Transfer, amount, currency, updatedAt)
	if err != nil {
		return nil, err
	}
	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) ListWithTotal(ctx context.Context, params ListTransactionsParams) (transactions []*Transaction, total int, err error) {
	if params.Sort == nil {
		sort := "DESC"
		params.Sort = &sort
	}
	if params.Limit == nil {
		limit := ListTransactionsLimit
		params.Limit = &limit
	}
	if params.Page == nil {
		page := 1
		params.Page = &page
	}

	return s.transactionRepo.ListWithTotalByAccountID(ctx, params)
}
