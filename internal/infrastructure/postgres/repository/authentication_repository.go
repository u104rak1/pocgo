package repository

import (
	"context"
	"database/sql"
	"errors"

	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type authenticationRepository struct {
	*Repository[model.Authentication]
}

func NewAuthenticationRepository(db *bun.DB) authDomain.IAuthenticationRepository {
	return &authenticationRepository{Repository: NewRepository[model.Authentication](db)}
}

func (r *authenticationRepository) Save(ctx context.Context, authentication *authDomain.Authentication) error {
	authModel := &model.Authentication{
		UserID:       authentication.UserIDString(),
		PasswordHash: authentication.PasswordHash(),
	}
	_, err := r.ExecDB(ctx).NewInsert().Model(authModel).On("CONFLICT (user_id) DO UPDATE").
		Set("password_hash = EXCLUDED.password_hash").
		Exec(ctx)
	return err
}

func (r *authenticationRepository) FindByUserID(ctx context.Context, userID idVO.UserID) (*authDomain.Authentication, error) {
	authModel := &model.Authentication{}
	if err := r.ExecDB(ctx).NewSelect().Model(authModel).Where("user_id = ?", userID.String()).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return authDomain.Reconstruct(userID.String(), authModel.PasswordHash)
}

func (r *authenticationRepository) ExistsByUserID(ctx context.Context, userID idVO.UserID) (bool, error) {
	exists, err := r.ExecDB(ctx).NewSelect().Model((*model.Authentication)(nil)).Where("user_id = ?", userID.String()).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}
