package user

import "context"

type IUserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByID(ctx context.Context, id UserID) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
