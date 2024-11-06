package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/ucho456job/pocgo/pkg/password"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/uptrace/bun"
)

const (
	JohnDoeWorkAccountID      = "01J9R8AJ1Q2YDH1X9836GS9D87"
	JohnDoePrivateAccountID   = "01J9R8AS3G2EA5723HB3E97QZE"
	JaneSmithWorkAccountID    = "01J9R8B042BTZEH5J9H5VR5TPM"
	JaneSmithPrivateAccountID = "01J9R8B83C89Q0JTAAB1YR1NHA"
)

func saveAccount(db *bun.DB) error {
	passwordHash, _ := password.Encode("1234")
	data := []model.Account{
		{
			ID:           JohnDoeWorkAccountID,
			UserID:       JohnDoeID,
			Name:         "work",
			PasswordHash: passwordHash,
			Balance:      100000,
			CurrencyID:   JPYID,
			UpdatedAt:    timer.Now(),
		},
		{
			ID:           JohnDoePrivateAccountID,
			UserID:       JohnDoeID,
			Name:         "private",
			PasswordHash: passwordHash,
			Balance:      200000,
			CurrencyID:   JPYID,
			UpdatedAt:    timer.Now(),
		},
		{
			ID:           JaneSmithWorkAccountID,
			UserID:       JaneSmithID,
			Name:         "work",
			PasswordHash: passwordHash,
			Balance:      3000.55,
			CurrencyID:   USDID,
			UpdatedAt:    timer.Now(),
		},
		{
			ID:           JaneSmithPrivateAccountID,
			UserID:       JaneSmithID,
			Name:         "private",
			PasswordHash: passwordHash,
			Balance:      4000.55,
			CurrencyID:   USDID,
			UpdatedAt:    timer.Now(),
		},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
