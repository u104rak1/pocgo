package signup_usecase

import (
	"context"
	"time"

	unitofwork "github.com/ucho456job/pocgo/internal/application/unit_of_work"
	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	authenticationDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCmd) (*SignupDTO, error)
}

type signupUsecase struct {
	accountRepo               accountDomain.IAccountRepository
	authenticationRepo        authenticationDomain.IAuthenticationRepository
	userRepo                  userDomain.IUserRepository
	accessTokenService        authenticationDomain.AccessTokenService
	verifyEmailUniquenessServ userDomain.VerifyEmailUniquenessService
	unitOfWork                unitofwork.IUnitOfWork
}

func NewSignupUsecase(
	accountRepo accountDomain.IAccountRepository,
	authenticationRepo authenticationDomain.IAuthenticationRepository,
	userRepo userDomain.IUserRepository,
	accessTokenService authenticationDomain.AccessTokenService,
	isEmailDuplicateServ userDomain.VerifyEmailUniquenessService,
	unitOfWork unitofwork.IUnitOfWork,
) ISignupUsecase {
	return &signupUsecase{
		accountRepo:               accountRepo,
		authenticationRepo:        authenticationRepo,
		userRepo:                  userRepo,
		accessTokenService:        accessTokenService,
		verifyEmailUniquenessServ: isEmailDuplicateServ,
		unitOfWork:                unitOfWork,
	}
}

func (u *signupUsecase) Run(ctx context.Context, cmd SignupCmd) (*SignupDTO, error) {
	var err error

	email := cmd.User.Email
	if err = u.verifyEmailUniquenessServ.Run(ctx, email); err != nil {
		return nil, err
	}

	userID := ulid.New()
	user, err := userDomain.New(userID, cmd.User.Name, email)
	if err != nil {
		return nil, err
	}

	accountID := ulid.New()
	account, err := accountDomain.New(
		accountID, userID, cmd.Account.Name, cmd.Account.Password,
		cmd.Account.Balance, cmd.Account.Currency, time.Now(),
	)
	if err != nil {
		return nil, err
	}

	authenticationID := ulid.New()
	authentication, err := authenticationDomain.New(authenticationID, userID, cmd.User.Password)
	if err != nil {
		return nil, err
	}

	err = u.unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
		if err = u.userRepo.Save(ctx, user); err != nil {
			return err
		}

		if err = u.accountRepo.Save(ctx, account); err != nil {
			return err
		}

		if err = u.authenticationRepo.Save(ctx, authentication); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := u.accessTokenService.Generate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return newSignupDTO(user, account, accessToken), nil
}
