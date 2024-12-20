package authentication

import (
	"context"

	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type IAuthenticationRepository interface {
	Save(ctx context.Context, authentication *Authentication) error
	FindByUserID(ctx context.Context, userID idVO.UserID) (*Authentication, error)
	ExistsByUserID(ctx context.Context, userID idVO.UserID) (bool, error)
}
