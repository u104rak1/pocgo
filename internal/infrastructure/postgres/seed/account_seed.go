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
	passwordHash, _ := password.Encode("1234")
	data := []model.AccountModel{
		{
			ID:           johnDoeWorkAccountID,
			UserID:       johnDoeID,
			Name:         "work",
			PasswordHash: passwordHash,
			Balance:      100000,
			CurrencyID:   jpyID,
			UpdatedAt:    time.Now(),
		},
		{
			ID:           johnDoePrivateAccountID,
			UserID:       johnDoeID,
			Name:         "private",
			PasswordHash: passwordHash,
			Balance:      200000,
			CurrencyID:   jpyID,
			UpdatedAt:    time.Now(),
		},
		{
			ID:           janeSmithWorkAccountID,
			UserID:       janeSmithID,
			Name:         "work",
			PasswordHash: passwordHash,
			Balance:      300000,
			CurrencyID:   jpyID,
			UpdatedAt:    time.Now(),
		},
		{
			ID:           janeSmithPrivateAccountID,
			UserID:       janeSmithID,
			Name:         "private",
			PasswordHash: passwordHash,
			Balance:      400000,
			CurrencyID:   jpyID,
			UpdatedAt:    time.Now(),
		},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
