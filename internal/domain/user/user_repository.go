package user

import (
	"context"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IUserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id idVO.UserID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByID(ctx context.Context, id idVO.UserID) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
