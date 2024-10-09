package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

const (
	deposit  = "DEPOSIT"
	withdraw = "WITHDRAW"
	transfer = "TRANSFER"
)

func SaveTransactionTypeMaster(db *bun.DB) error {
	data := []model.TransactionTypeMasterModel{
		{Type: deposit},
		{Type: withdraw},
		{Type: transfer},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
