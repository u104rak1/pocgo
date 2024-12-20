package account

import (
	"context"

	unitofwork "github.com/u104rak1/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
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
	userID, err := idVO.UserIDFromString(cmd.UserID)
	if err != nil {
		return nil, err
	}

	balance := 0.0
	account, err := accountDomain.New(
		userID, balance, cmd.Name, cmd.Password, cmd.Currency,
	)
	if err != nil {
		return nil, err
	}

	err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.userServ.EnsureUserExists(ctx, userID); err != nil {
			return err
		}

		if err := u.accountServ.CheckLimit(ctx, userID); err != nil {
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
		ID:        account.IDString(),
		UserID:    account.UserIDString(),
		Name:      account.Name(),
		Balance:   account.Balance().Amount(),
		Currency:  account.Balance().Currency(),
		UpdatedAt: account.UpdatedAtString(),
	}, nil
}
