package user

import (
	"context"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IReadUserUsecase interface {
	Run(ctx context.Context, cmd ReadUserCommand) (*ReadUserDTO, error)
}

type readUserUsecase struct {
	userServ userDomain.IUserService
}

func NewReadUserUsecase(userService userDomain.IUserService) IReadUserUsecase {
	return &readUserUsecase{
		userServ: userService,
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
	userID, err := idVO.UserIDFromString(cmd.ID)
	if err != nil {
		return nil, err
	}

	user, err := u.userServ.FindUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ReadUserDTO{
		ID:    user.IDString(),
		Name:  user.Name(),
		Email: user.Email(),
	}, nil
}
