package authentication

import (
	"context"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
)

type IAuthenticationRepository interface {
	Save(ctx context.Context, authentication *Authentication) error
	FindByUserID(ctx context.Context, userID userDomain.UserID) (*Authentication, error)
	ExistsByUserID(ctx context.Context, userID userDomain.UserID) (bool, error)
}
