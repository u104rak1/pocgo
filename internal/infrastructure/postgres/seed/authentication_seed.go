package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/ucho456job/pocgo/pkg/password"
	"github.com/uptrace/bun"
)

func saveAuthentication(db *bun.DB) error {
	passwordHash, _ := password.Encode("password")
	data := []model.Authentication{
		{UserID: JohnDoeID, PasswordHash: passwordHash},
		{UserID: JaneSmithID, PasswordHash: passwordHash},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
