package transaction

import (
	"context"

	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ITransactionService interface {
	Deposit(ctx context.Context, account *accountDomain.Account, amount float64, currency string) (*Transaction, error)
	Withdraw(ctx context.Context, account *accountDomain.Account, amount float64, currency string) (*Transaction, error)
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

	transactionID := ulid.New()
	transaction, err := New(transactionID, account.ID(), nil, Deposit, amount, currency, updatedAt)
	if err != nil {
		return nil, err
	}
	if err := s.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) Withdraw(
	ctx context.Context,
	account *accountDomain.Account,
	amount float64,
	currency string,
) (*Transaction, error) {
	if err := account.Withdraw(amount, currency); err != nil {
		return nil, err
	}
	updatedAt := timer.Now()
	account.ChangeUpdatedAt(updatedAt)
	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	transactionID := ulid.New()
	transaction, err := New(transactionID, account.ID(), nil, Withdraw, amount, currency, updatedAt)
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
	if err := senderAccount.Withdraw(amount, currency); err != nil {
		return nil, err
	}
	if err := receiverAccount.Deposit(amount, currency); err != nil {
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

	transactionID := ulid.New()
	receiverAccountID := receiverAccount.ID()
	transaction, err := New(transactionID, senderAccount.ID(), &receiverAccountID, Transfer, amount, currency, updatedAt)
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
