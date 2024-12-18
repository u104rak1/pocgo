package repository

import (
	"context"
	"database/sql"
	"errors"

	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/uptrace/bun"
)

type userRepository struct {
	*Repository[model.User]
}

func NewUserRepository(db *bun.DB) userDomain.IUserRepository {
	return &userRepository{Repository: NewRepository[model.User](db)}
}

func (r *userRepository) Save(ctx context.Context, user *userDomain.User) error {
	userModel := &model.User{
		ID:    user.ID(),
		Email: user.Email(),
		Name:  user.Name(),
	}
	_, err := r.execDB(ctx).NewInsert().Model(userModel).On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("email = EXCLUDED.email").
		Exec(ctx)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*userDomain.User, error) {
	userModel := &model.User{}
	if err := r.execDB(ctx).NewSelect().Model(userModel).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return userDomain.New(userModel.ID, userModel.Name, userModel.Email)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	userModel := &model.User{}
	if err := r.execDB(ctx).NewSelect().Model(userModel).Where("email = ?", email).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return userDomain.New(userModel.ID, userModel.Name, userModel.Email)
}

func (r *userRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	return r.execDB(ctx).NewSelect().Model((*model.User)(nil)).Where("id = ?", id).Exists(ctx)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.execDB(ctx).NewSelect().Model((*model.User)(nil)).Where("email = ?", email).Exists(ctx)
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.execDB(ctx).NewUpdate().
		Model(&model.User{ID: id, DeletedAt: timer.Now()}).
		Column("deleted_at").
		WherePK().
		Exec(ctx)
	return err
}
