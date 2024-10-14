package signup

import (
	"context"

	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error)
}

type signupUsecase struct {
	createUserUsecase    userApp.ICreateUserUsecase
	createAccountUsecase accountApp.ICreateAccountUsecase
	accessTokenService   authDomain.AccessTokenService
	unitOfWork           unitofwork.IUnitOfWorkWithResult[SignupDTO]
}

func NewSignupUsecase(
	createUserUsecase userApp.ICreateUserUsecase,
	createAccountUsecase accountApp.ICreateAccountUsecase,
	accessTokenService authDomain.AccessTokenService,
	unitOfWork unitofwork.IUnitOfWorkWithResult[SignupDTO],
) ISignupUsecase {
	return &signupUsecase{
		createUserUsecase:    createUserUsecase,
		createAccountUsecase: createAccountUsecase,
		accessTokenService:   accessTokenService,
		unitOfWork:           unitOfWork,
	}
}

type SignupCommand struct {
	User    userApp.CreateUserCommand
	Account accountApp.CreateAccountCommand
}

type SignupDTO struct {
	User        userApp.CreateUserDTO
	Account     accountApp.CreateAccountDTO
	AccessToken string
}

func (u *signupUsecase) Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error) {
	dto, err := u.unitOfWork.RunInTx(ctx, func(ctx context.Context, tx unitofwork.ITransaction) (*SignupDTO, error) {
		user, err := u.createUserUsecase.Run(ctx, cmd.User)
		if err != nil {
			return nil, err
		}

		cmd.Account.UserID = user.ID
		account, err := u.createAccountUsecase.Run(ctx, cmd.Account)
		if err != nil {
			return nil, err
		}

		return &SignupDTO{
			User:    *user,
			Account: *account,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := u.accessTokenService.Generate(ctx, dto.User.ID)
	if err != nil {
		return nil, err
	}
	dto.AccessToken = accessToken

	return dto, nil
}
