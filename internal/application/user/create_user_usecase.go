package user

import (
	"context"

	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ICreateUserUsecase interface {
	Run(ctx context.Context, cmd CreateUserCommand) (*CreateUserDTO, error)
}

type createUserUsecase struct {
	userRepo                           userDomain.IUserRepository
	userServ                           userDomain.IUserService
	authRepo                           authDomain.IAuthenticationRepository
	verifyAuthenticationUniquenessServ authDomain.IVerifyAuthenticationUniquenessService
}

func NewCreateUserUsecase(
	userRepo userDomain.IUserRepository,
	authRepo authDomain.IAuthenticationRepository,
	userServ userDomain.IUserService,
	verifyAuthenticationUniquenessServ authDomain.IVerifyAuthenticationUniquenessService,
) ICreateUserUsecase {
	return &createUserUsecase{
		userRepo:                           userRepo,
		authRepo:                           authRepo,
		userServ:                           userServ,
		verifyAuthenticationUniquenessServ: verifyAuthenticationUniquenessServ,
	}
}

type CreateUserCommand struct {
	Name     string
	Email    string
	Password string
}

type CreateUserDTO struct {
	ID    string
	Name  string
	Email string
}

func (u *createUserUsecase) Run(ctx context.Context, cmd CreateUserCommand) (*CreateUserDTO, error) {
	userID := ulid.New()
	if err := u.createUser(ctx, userID, cmd); err != nil {
		return nil, err
	}

	if err := u.createAuthentication(ctx, userID, cmd); err != nil {
		return nil, err
	}

	return &CreateUserDTO{
		ID:    userID,
		Name:  cmd.Name,
		Email: cmd.Email,
	}, nil
}

func (u *createUserUsecase) createUser(ctx context.Context, userID string, cmd CreateUserCommand) (err error) {
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

func (u *createUserUsecase) createAuthentication(ctx context.Context, userID string, cmd CreateUserCommand) (err error) {
	if err = u.verifyAuthenticationUniquenessServ.Run(ctx, userID); err != nil {
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
