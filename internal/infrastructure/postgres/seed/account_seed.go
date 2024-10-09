package seed

import (
	"context"
	"time"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/ucho456job/pocgo/pkg/password"
	"github.com/uptrace/bun"
)

const (
	johnDoeWorkAccountID      = "01J9R8AJ1Q2YDH1X9836GS9D87"
	johnDoePrivateAccountID   = "01J9R8AS3G2EA5723HB3E97QZE"
	janeSmithWorkAccountID    = "01J9R8B042BTZEH5J9H5VR5TPM"
	janeSmithPrivateAccountID = "01J9R8B83C89Q0JTAAB1YR1NHA"
)

func saveAccount(db *bun.DB) error {
	data := []model.AccountModel{
		{
			ID:            johnDoeWorkAccountID,
			UserID:        johnDoeID,
			Name:          "work",
			PasswordHash:  password.Encode("1234"),
			Balance:       100000,
			CurrencyID:    jpyID,
			LastUpdatedAt: time.Now(),
		},
		{
			ID:            johnDoePrivateAccountID,
			UserID:        johnDoeID,
			Name:          "private",
			PasswordHash:  password.Encode("1234"),
			Balance:       200000,
			CurrencyID:    jpyID,
			LastUpdatedAt: time.Now(),
		},
		{
			ID:            janeSmithWorkAccountID,
			UserID:        janeSmithID,
			Name:          "work",
			PasswordHash:  password.Encode("1234"),
			Balance:       300000,
			CurrencyID:    jpyID,
			LastUpdatedAt: time.Now(),
		},
		{
			ID:            janeSmithPrivateAccountID,
			UserID:        janeSmithID,
			Name:          "private",
			PasswordHash:  password.Encode("1234"),
			Balance:       400000,
			CurrencyID:    jpyID,
			LastUpdatedAt: time.Now(),
		},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
