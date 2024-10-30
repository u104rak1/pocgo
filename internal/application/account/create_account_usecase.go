package account

import (
	"context"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ICreateAccountUsecase interface {
	Run(ctx context.Context, cmd CreateAccountCommand) (*CreateAccountDTO, error)
}

type createAccountUsecase struct {
	accountRepo accountDomain.IAccountRepository
	accountServ accountDomain.IAccountService
	userServ    userDomain.IUserService
	unitOfWork  unitofwork.IUnitOfWork
}

func NewCreateAccountUsecase(
	accountRepository accountDomain.IAccountRepository,
	accountService accountDomain.IAccountService,
	userService userDomain.IUserService,
	unitOfWork unitofwork.IUnitOfWork,
) ICreateAccountUsecase {
	return &createAccountUsecase{
		accountRepo: accountRepository,
		accountServ: accountService,
		userServ:    userService,
		unitOfWork:  unitOfWork,
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

	err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.userServ.EnsureUserExists(ctx, cmd.UserID); err != nil {
			return err
		}

		if err := u.accountServ.CheckLimit(ctx, cmd.UserID); err != nil {
			return err
		}

		if err = u.accountRepo.Save(ctx, account); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
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
