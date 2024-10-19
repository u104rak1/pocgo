package authentication

// mockgen -source=internal/domain/authentication/authentication_repository.go -destination=internal/domain/authentication/mock/mock_authentication_repository.go -package=mock

import "context"

type IAuthenticationRepository interface {
	Save(ctx context.Context, authentication *Authentication) error
	FindByUserID(ctx context.Context, userID string) (*Authentication, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}
