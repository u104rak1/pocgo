package user_repository

import (
	"context"

	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type IUserRepository interface {
	Save(ctx context.Context, user *userDomain.User) error
	FindByID(ctx context.Context, id string) (*userDomain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *userDomain.User) error {
	userModel := r.factoryUserModel(user)
	_, err := r.db.NewInsert().Model(userModel).On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("email = EXCLUDED.email").
		Exec(ctx)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*userDomain.User, error) {
	userModel := new(model.UserModel)
	err := r.db.NewSelect().Model(userModel).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return userDomain.New(userModel.ID, userModel.Email, userModel.Name)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().Model((*model.UserModel)(nil)).Where("email = ?", email).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) factoryUserModel(user *userDomain.User) *model.UserModel {
	return &model.UserModel{
		ID:    user.ID(),
		Email: user.Email(),
		Name:  user.Name(),
	}
}
