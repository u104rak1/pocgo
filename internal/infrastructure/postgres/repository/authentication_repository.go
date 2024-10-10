package repository

import (
	"context"

	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type IAuthenticationRepository interface {
	Save(ctx context.Context, authentication *authDomain.Authentication) error
	FindByUserID(ctx context.Context, userID string) (*authDomain.Authentication, error)
}

type authenticationRepository struct {
	db *bun.DB
}

func NewAuthenticationRepository(db *bun.DB) IAuthenticationRepository {
	return &authenticationRepository{db: db}
}

func (r *authenticationRepository) Save(ctx context.Context, authentication *authDomain.Authentication) error {
	authModel := &model.Authentication{
		UserID:       authentication.UserID(),
		PasswordHash: authentication.PasswordHash(),
	}
	_, err := r.db.NewInsert().Model(authModel).On("CONFLICT (user_id) DO UPDATE").
		Set("password_hash = EXCLUDED.password_hash").
		Exec(ctx)
	return err
}

func (r *authenticationRepository) FindByUserID(ctx context.Context, userID string) (*authDomain.Authentication, error) {
	authModel := &model.Authentication{}
	if err := r.db.NewSelect().Model(authModel).Where("user_id = ?", userID).Scan(ctx); err != nil {
		return nil, err
	}
	return authDomain.Reconstruct(userID, authModel.PasswordHash)
}
