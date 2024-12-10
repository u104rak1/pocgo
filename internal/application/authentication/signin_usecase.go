package authentication

import (
	"context"

	"github.com/u104rak1/pocgo/internal/config"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
)

type ISigninUsecase interface {
	Run(ctx context.Context, cmd SigninCommand) (*SigninDTO, error)
}

type signinUsecase struct {
	authServ authDomain.IAuthenticationService
}

func NewSigninUsecase(
	authenticationService authDomain.IAuthenticationService,
) ISigninUsecase {
	return &signinUsecase{
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

	env := config.NewEnv()
	token, err := u.authServ.GenerateAccessToken(ctx, userID, []byte(env.JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return &SigninDTO{
		AccessToken: token,
	}, nil
}
