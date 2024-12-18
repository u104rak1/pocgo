package user

import (
	"context"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
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
	//TODO: エラーメッセージがユースケースに出てくる場合はdomain serviceに書く
	user, err := u.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, userDomain.ErrNotFound
	}

	return &ReadUserDTO{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}, nil
}
