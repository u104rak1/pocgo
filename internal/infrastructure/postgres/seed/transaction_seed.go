package seed

import (
	"context"
	"time"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func saveTransaction(db *bun.DB) error {
	transactions := []model.TransactionModel{
		{
			ID:                "01J9RFMD0GQ3Q36RP34HBSBYHM",
			AccountID:         johnDoeWorkAccountID,
			ReceiverAccountID: nil,
			Type:              deposit,
			Amount:            100000,
			CurrencyID:        jpyID,
			TransactionAt:     time.Now(),
		},
		{
			ID:                "01J9RFS63XVCFBC4ND478A9FWB",
			AccountID:         johnDoePrivateAccountID,
			ReceiverAccountID: nil,
			Type:              deposit,
			Amount:            200000,
			CurrencyID:        jpyID,
			TransactionAt:     time.Now(),
		},
		{
			ID:                "01J9RFTNH5J6W4XNBB7G8A37HC",
			AccountID:         janeSmithWorkAccountID,
			ReceiverAccountID: nil,
			Type:              deposit,
			Amount:            300000,
			CurrencyID:        jpyID,
			TransactionAt:     time.Now(),
		},
		{
			ID:                "01J9RFVNM10Y0CMC26XG85SB7S",
			AccountID:         janeSmithPrivateAccountID,
			ReceiverAccountID: nil,
			Type:              deposit,
			Amount:            400000,
			CurrencyID:        jpyID,
			TransactionAt:     time.Now(),
		},
	}

	if _, err := db.NewInsert().Model(&transactions).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}
