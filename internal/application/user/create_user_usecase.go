package user_usecase

import (
	"context"

	authenticationDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

type ICreateUserUsecase interface {
	Run(ctx context.Context, cmd CreateUserCommand) (*CreateUserDTO, error)
}

type createUserUsecase struct {
	userRepo                           userDomain.IUserRepository
	authenticationRepo                 authenticationDomain.IAuthenticationRepository
	verifyEmailUniquenessServ          userDomain.VerifyEmailUniquenessService
	verifyAuthenticationUniquenessServ authenticationDomain.VerifyAuthenticationUniquenessService
}

func NewCreateUserUsecase(
	userRepo userDomain.IUserRepository,
	authenticationRepo authenticationDomain.IAuthenticationRepository,
	verifyEmailUniquenessServ userDomain.VerifyEmailUniquenessService,
	verifyAuthenticationUniquenessServ authenticationDomain.VerifyAuthenticationUniquenessService,
) ICreateUserUsecase {
	return &createUserUsecase{
		userRepo:                           userRepo,
		authenticationRepo:                 authenticationRepo,
		verifyEmailUniquenessServ:          verifyEmailUniquenessServ,
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
	var err error

	userID := ulid.New()
	if err = u.createUser(ctx, userID, cmd); err != nil {
		return nil, err
	}

	if err = u.createAuthentication(ctx, userID, cmd); err != nil {
		return nil, err
	}

	return &CreateUserDTO{
		ID:    userID,
		Name:  cmd.Name,
		Email: cmd.Email,
	}, nil
}

func (u *createUserUsecase) createUser(ctx context.Context, userID string, cmd CreateUserCommand) (err error) {
	if err = u.verifyEmailUniquenessServ.Run(ctx, cmd.Email); err != nil {
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

	authentication, err := authenticationDomain.New(userID, cmd.Password)
	if err != nil {
		return err
	}

	if err = u.authenticationRepo.Save(ctx, authentication); err != nil {
		return err
	}

	return nil
}
