package seed

import (
	"context"

	transactionDomain "github.com/u104raki/pocgo/internal/domain/transaction"
	"github.com/u104raki/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104raki/pocgo/pkg/timer"
	"github.com/uptrace/bun"
)

func saveTransaction(db *bun.DB) error {
	transactions := []model.Transaction{
		{
			ID:                "01J9RFMD0GQ3Q36RP34HBSBYHM",
			AccountID:         JohnDoeWorkAccountID,
			ReceiverAccountID: nil,
			OperationType:     transactionDomain.Deposit,
			Amount:            100000,
			CurrencyID:        JPYID,
			TransactionAt:     timer.Now(),
		},
		{
			ID:                "01J9RFS63XVCFBC4ND478A9FWB",
			AccountID:         JohnDoePrivateAccountID,
			ReceiverAccountID: nil,
			OperationType:     transactionDomain.Deposit,
			Amount:            200000,
			CurrencyID:        JPYID,
			TransactionAt:     timer.Now(),
		},
		{
			ID:                "01J9RFTNH5J6W4XNBB7G8A37HC",
			AccountID:         JaneSmithWorkAccountID,
			ReceiverAccountID: nil,
			OperationType:     transactionDomain.Deposit,
			Amount:            3000.55,
			CurrencyID:        USDID,
			TransactionAt:     timer.Now(),
		},
		{
			ID:                "01J9RFVNM10Y0CMC26XG85SB7S",
			AccountID:         JaneSmithPrivateAccountID,
			ReceiverAccountID: nil,
			OperationType:     transactionDomain.Deposit,
			Amount:            4000.55,
			CurrencyID:        USDID,
			TransactionAt:     timer.Now(),
		},
	}

	if _, err := db.NewInsert().Model(&transactions).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}
