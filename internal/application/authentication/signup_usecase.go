package authentication

import (
	"context"

	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/environment"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error)
}

type signupUsecase struct {
	createUserUC    userApp.ICreateUserUsecase
	createAccountUC accountApp.ICreateAccountUsecase
	authServ        authDomain.IAuthenticationService
	unitOfWork      unitofwork.IUnitOfWorkWithResult[SignupDTO]
}

func NewSignupUsecase(
	createUserUsecase userApp.ICreateUserUsecase,
	createAccountUsecase accountApp.ICreateAccountUsecase,
	authService authDomain.IAuthenticationService,
	unitOfWork unitofwork.IUnitOfWorkWithResult[SignupDTO],
) ISignupUsecase {
	return &signupUsecase{
		createUserUC:    createUserUsecase,
		createAccountUC: createAccountUsecase,
		authServ:        authService,
		unitOfWork:      unitOfWork,
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
	dto, err := u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*SignupDTO, error) {
		user, err := u.createUserUC.Run(ctx, cmd.User)
		if err != nil {
			return nil, err
		}

		cmd.Account.UserID = user.ID
		account, err := u.createAccountUC.Run(ctx, cmd.Account)
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

	env := environment.New()
	accessToken, err := u.authServ.GenerateAccessToken(ctx, dto.User.ID, []byte(env.JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}
	dto.AccessToken = accessToken

	return dto, nil
}
