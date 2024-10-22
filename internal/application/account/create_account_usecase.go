package account

import (
	"context"

	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ICreateAccountUsecase interface {
	Run(ctx context.Context, cmd CreateAccountCommand) (*CreateAccountDTO, error)
}

type createAccountUsecase struct {
	accountRepo accountDomain.IAccountRepository
}

func NewCreateAccountUsecase(accountRepo accountDomain.IAccountRepository) ICreateAccountUsecase {
	return &createAccountUsecase{
		accountRepo: accountRepo,
	}
}

type CreateAccountCommand struct {
	UserID   string
	Name     string
	Password string
	Currency string
}

type CreateAccountDTO struct {
	ID        string
	UserID    string
	Name      string
	Balance   float64
	Currency  string
	UpdatedAt string
}

func (u *createAccountUsecase) Run(ctx context.Context, cmd CreateAccountCommand) (*CreateAccountDTO, error) {
	accountID := ulid.New()
	account, err := accountDomain.New(
		accountID, cmd.UserID, cmd.Name, cmd.Password,
		0, cmd.Currency, timer.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err = u.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	return &CreateAccountDTO{
		ID:        account.ID(),
		UserID:    account.UserID(),
		Name:      account.Name(),
		Balance:   account.Balance().Amount(),
		Currency:  account.Balance().Currency(),
		UpdatedAt: account.UpdatedAt().String(),
	}, nil
}
