package authentication

import (
	"context"

	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
)

type ISigninUsecase interface {
	Run(ctx context.Context, cmd SigninCommand) (*SigninDTO, error)
}

type signinUsecase struct {
	authServ authDomain.IAuthenticationService
	jwtServ  IJWTService
}

func NewSigninUsecase(
	authenticationService authDomain.IAuthenticationService,
	jwtService IJWTService,
) ISigninUsecase {
	return &signinUsecase{
		authServ: authenticationService,
		jwtServ:  jwtService,
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

	token, err := u.jwtServ.GenerateAccessToken(userID.String())
	if err != nil {
		return nil, err
	}

	return &SigninDTO{
		AccessToken: token,
	}, nil
}
