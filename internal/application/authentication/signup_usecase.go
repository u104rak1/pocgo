package authentication

import (
	"context"

	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type ISignupUsecase interface {
	Run(ctx context.Context, cmd SignupCommand) (*SignupDTO, error)
}

type signupUsecase struct {
	userRepo userDomain.IUserRepository
	userServ userDomain.IUserService
	authRepo authDomain.IAuthenticationRepository
	authServ authDomain.IAuthenticationService
	jwtServ  IJWTService
}

func NewSignupUsecase(
	userRepository userDomain.IUserRepository,
	authRepository authDomain.IAuthenticationRepository,
	userService userDomain.IUserService,
	authService authDomain.IAuthenticationService,
	jwtService IJWTService,
) ISignupUsecase {
	return &signupUsecase{
		userRepo: userRepository,
		authRepo: authRepository,
		userServ: userService,
		authServ: authService,
		jwtServ:  jwtService,
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
	userID, err := u.createUser(ctx, cmd)
	if err != nil {
		return nil, err
	}

	if err := u.createAuthentication(ctx, userID, cmd); err != nil {
		return nil, err
	}

	accessToken, err := u.jwtServ.GenerateAccessToken(userID.String())
	if err != nil {
		return nil, err
	}

	return &SignupDTO{
		User: SignupUserDTO{
			ID:    userID.String(),
			Name:  cmd.Name,
			Email: cmd.Email,
		},
		AccessToken: accessToken,
	}, nil
}

func (u *signupUsecase) createUser(ctx context.Context, cmd SignupCommand) (*idVO.UserID, error) {
	if err := u.userServ.VerifyEmailUniqueness(ctx, cmd.Email); err != nil {
		return nil, err
	}

	user, err := userDomain.New(cmd.Name, cmd.Email)
	if err != nil {
		return nil, err
	}

	if err = u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	userID := user.ID()
	return &userID, nil
}

func (u *signupUsecase) createAuthentication(ctx context.Context, userID *idVO.UserID, cmd SignupCommand) error {
	if err := u.authServ.VerifyUniqueness(ctx, *userID); err != nil {
		return err
	}

	authentication, err := authDomain.New(*userID, cmd.Password)
	if err != nil {
		return err
	}

	if err = u.authRepo.Save(ctx, authentication); err != nil {
		return err
	}

	return nil
}
