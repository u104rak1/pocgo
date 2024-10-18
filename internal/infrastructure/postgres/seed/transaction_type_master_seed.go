package seed

import (
	"context"

	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func saveTransactionTypeMaster(db *bun.DB) error {
	data := []model.TransactionTypeMaster{
		{Type: transactionDomain.Deposit},
		{Type: transactionDomain.Withdraw},
		{Type: transactionDomain.Transfer},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
