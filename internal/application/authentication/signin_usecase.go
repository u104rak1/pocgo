package authentication

import (
	"context"

	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/environment"
)

type ISigninUsecase interface {
	Run(ctx context.Context, cmd SigninCommand) (*SigninDTO, error)
}

type signinUsecase struct {
	userRepo userDomain.IUserRepository
	authRepo authDomain.IAuthenticationRepository
	authServ authDomain.IAuthenticationService
}

func NewSigninUsecase(
	userRepository userDomain.IUserRepository,
	authenticationRepository authDomain.IAuthenticationRepository,
	authenticationService authDomain.IAuthenticationService,
) ISigninUsecase {
	return &signinUsecase{
		userRepo: userRepository,
		authRepo: authenticationRepository,
		authServ: authenticationService,
	}
}

type SigninCommand struct {
	Email    string
	Password string
}

type SigninDTO struct {
	AccessToken string
}

func (u *signinUsecase) Run(ctx context.Context, cmd SigninCommand) (*SigninDTO, error) {
	userID, err := u.authServ.Authenticate(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	env := environment.New()
	token, err := u.authServ.GenerateAccessToken(ctx, userID, []byte(env.JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return &SigninDTO{
		AccessToken: token,
	}, nil
}
