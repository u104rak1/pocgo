package signup_usecase

import (
	"context"

	accountUC "github.com/ucho456job/pocgo/internal/application/account"
	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	userUC "github.com/ucho456job/pocgo/internal/application/user"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error)
}

type signupUsecase struct {
	createUserUC    userUC.ICreateUserUsecase
	createAccountUC accountUC.ICreateAccountUsecase
	accessTokenServ authDomain.AccessTokenService
	unitOfWork      unitofwork.IUnitOfWorkWithResult[*SignupDTO]
}

func NewSignupUsecase(
	createUserUC userUC.ICreateUserUsecase,
	createAccountUC accountUC.ICreateAccountUsecase,
	accessTokenServ authDomain.AccessTokenService,
	unitOfWork unitofwork.IUnitOfWorkWithResult[*SignupDTO],
) ISignupUsecase {
	return &signupUsecase{
		createUserUC:    createUserUC,
		createAccountUC: createAccountUC,
		accessTokenServ: accessTokenServ,
		unitOfWork:      unitOfWork,
	}
}

type SignupCommand struct {
	User    userUC.CreateUserCommand
	Account accountUC.CreateAccountCommand
}

type SignupDTO struct {
	User        userUC.CreateUserDTO
	Account     accountUC.CreateAccountDTO
	AccessToken string
}

func (u *signupUsecase) Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error) {
	dto, err := u.unitOfWork.RunInTx(ctx, func(ctx context.Context) (*SignupDTO, error) {
		user, err := u.createUserUC.Run(ctx, cmd.User)
		if err != nil {
			return nil, err
		}

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

	accessToken, err := u.accessTokenServ.Generate(ctx, dto.User.ID)
	if err != nil {
		return nil, err
	}
	dto.AccessToken = accessToken

	return dto, nil
}
