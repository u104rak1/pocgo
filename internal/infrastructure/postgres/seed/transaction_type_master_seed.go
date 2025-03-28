package seed

import (
	"context"

	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func saveOperationTypeMaster(db *bun.DB) error {
	data := []model.OperationTypeMaster{
		{Type: transactionDomain.Deposit},
		{Type: transactionDomain.Withdrawal},
		{Type: transactionDomain.Transfer},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
