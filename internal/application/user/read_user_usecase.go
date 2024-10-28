package user

import (
	"context"

	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
)

type IReadUserUsecase interface {
	Run(ctx context.Context, cmd ReadUserCommand) (*ReadUserDTO, error)
}

type readUserUsecase struct {
	userRepo userDomain.IUserRepository
}

func NewReadUserUsecase(userRepository userDomain.IUserRepository) IReadUserUsecase {
	return &readUserUsecase{
		userRepo: userRepository,
	}
}

type ReadUserCommand struct {
	ID string
}

type ReadUserDTO struct {
	ID    string
	Name  string
	Email string
}

func (u *readUserUsecase) Run(ctx context.Context, cmd ReadUserCommand) (*ReadUserDTO, error) {
	user, err := u.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	return &ReadUserDTO{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}, nil
}
