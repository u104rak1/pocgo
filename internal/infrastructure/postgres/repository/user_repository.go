package repository

import (
	"context"
	"time"

	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) userDomain.IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *userDomain.User) error {
	tx := getTx(ctx)

	var execDB bun.IDB = r.db
	if tx != nil {
		execDB = tx
	}

	userModel := &model.User{
		ID:    user.ID(),
		Email: user.Email(),
		Name:  user.Name(),
	}
	_, err := execDB.NewInsert().Model(userModel).On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("email = EXCLUDED.email").
		Exec(ctx)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*userDomain.User, error) {
	userModel := &model.User{}
	if err := r.db.NewSelect().Model(userModel).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return userDomain.New(userModel.ID, userModel.Email, userModel.Name)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().Model((*model.User)(nil)).Where("email = ?", email).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	tx := getTx(ctx)

	var execDB bun.IDB = r.db
	if tx != nil {
		execDB = tx
	}

	_, err := execDB.NewUpdate().
		Model(&model.User{ID: id, DeletedAt: time.Now()}).
		Column("deleted_at").
		WherePK().
		Exec(ctx)
	return err
}
