package authentication

import (
	"context"

	"github.com/u104rak1/pocgo/internal/config"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error)
}

type signupUsecase struct {
	userRepo userDomain.IUserRepository
	userServ userDomain.IUserService
	authRepo authDomain.IAuthenticationRepository
	authServ authDomain.IAuthenticationService
}

func NewSignupUsecase(
	userRepository userDomain.IUserRepository,
	authRepository authDomain.IAuthenticationRepository,
	userService userDomain.IUserService,
	authService authDomain.IAuthenticationService,
) ISignupUsecase {
	return &signupUsecase{
		userRepo: userRepository,
		authRepo: authRepository,
		userServ: userService,
		authServ: authService,
	}
}

type SignupCommand struct {
	Name     string
	Email    string
	Password string
}

type SignupDTO struct {
	User        SignupUserDTO
	AccessToken string
}

type SignupUserDTO struct {
	ID    string
	Name  string
	Email string
}

func (u *signupUsecase) Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error) {
	userID := ulid.New()
	if err := u.createUser(ctx, userID, cmd); err != nil {
		return nil, err
	}

	if err := u.createAuthentication(ctx, userID, cmd); err != nil {
		return nil, err
	}

	env := config.NewEnv()
	accessToken, err := u.authServ.GenerateAccessToken(ctx, userID, []byte(env.JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return &SignupDTO{
		User: SignupUserDTO{
			ID:    userID,
			Name:  cmd.Name,
			Email: cmd.Email,
		},
		AccessToken: accessToken,
	}, nil
}

func (u *signupUsecase) createUser(ctx context.Context, userID string, cmd SignupCommand) (err error) {
	if err = u.userServ.VerifyEmailUniqueness(ctx, cmd.Email); err != nil {
		return err
	}

	user, err := userDomain.New(userID, cmd.Name, cmd.Email)
	if err != nil {
		return err
	}

	if err = u.userRepo.Save(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *signupUsecase) createAuthentication(ctx context.Context, userID string, cmd SignupCommand) (err error) {
	if err = u.authServ.VerifyUniqueness(ctx, userID); err != nil {
		return err
	}

	authentication, err := authDomain.New(userID, cmd.Password)
	if err != nil {
		return err
	}

	if err = u.authRepo.Save(ctx, authentication); err != nil {
		return err
	}

	return nil
}
