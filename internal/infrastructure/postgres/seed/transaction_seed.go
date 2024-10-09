package seed

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func generateULIDWithTime(timestamp time.Time) string {
	entropy := ulid.Monotonic(nil, 0)
	return ulid.MustNew(ulid.Timestamp(timestamp), entropy).String()
}

func SaveTransactions(db *bun.DB) error {
	transactions := []model.TransactionModel{}

	// 入金・引き出しの取引
	for _, trans := range []struct {
		accountID string
		amount    float64
		transType string
	}{
		{johnDoeWorkAccountID, 10000, deposit},
		{johnDoeWorkAccountID, 5000, withdraw},
		{johnDoePrivateAccountID, 20000, deposit},
		{johnDoePrivateAccountID, 10000, withdraw},
		{janeSmithWorkAccountID, 15000, deposit},
		{janeSmithWorkAccountID, 7000, withdraw},
		{janeSmithPrivateAccountID, 25000, deposit},
		{janeSmithPrivateAccountID, 12000, withdraw},
	} {
		timestamp := time.Now()
		ulidID := generateULIDWithTime(timestamp)
		transactions = append(transactions, model.TransactionModel{
			ID:            ulidID,
			AccountID:     trans.accountID,
			Type:          trans.transType,
			Amount:        trans.amount,
			CurrencyID:    jpyID,
			TransactionAt: timestamp,
		})
	}

	// 送金の取引
	for _, trans := range []struct {
		accountID         string
		receiverAccountID string
		amount            float64
	}{
		{johnDoeWorkAccountID, johnDoePrivateAccountID, 3000},
		{johnDoePrivateAccountID, janeSmithWorkAccountID, 4000},
		{janeSmithWorkAccountID, janeSmithPrivateAccountID, 2000},
	} {
		timestamp := time.Now()
		ulidID := generateULIDWithTime(timestamp)
		receiverAccount := trans.receiverAccountID
		transactions = append(transactions, model.TransactionModel{
			ID:                ulidID,
			AccountID:         trans.accountID,
			ReceiverAccountID: &receiverAccount,
			Type:              transfer,
			Amount:            trans.amount,
			CurrencyID:        jpyID,
			TransactionAt:     timestamp,
		})
	}

	if _, err := db.NewInsert().Model(&transactions).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}
