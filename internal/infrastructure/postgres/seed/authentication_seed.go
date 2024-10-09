package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/ucho456job/pocgo/pkg/password"
	"github.com/uptrace/bun"
)

func saveAuthentication(db *bun.DB) error {
	data := []model.AuthenticationModel{
		{UserID: johnDoeID, PasswordHash: password.Encode("password")},
		{UserID: janeSmithID, PasswordHash: password.Encode("password")},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
