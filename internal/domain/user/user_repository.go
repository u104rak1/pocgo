package user

// mockgen -source=internal/domain/user/user_repository.go -destination=internal/domain/user/mock/mock_user_repository.go -package=mock

import "context"

type IUserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Delete(ctx context.Context, id string) error
}
