package authentication

import "context"

type IAuthenticationRepository interface {
	Save(ctx context.Context, authentication *Authentication) error
	FindByUserID(ctx context.Context, userID string) (*Authentication, error)
}
