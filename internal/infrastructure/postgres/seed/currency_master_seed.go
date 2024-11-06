package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

const (
	JPYID = "01J9R7YPV1FH1V0PPKVSB5C9TQ"
	USDID = "01J9R7ZQZQZQZQZQZQZQZQZQZQ"
)

func saveCurrencyMaster(db *bun.DB) error {
	data := []model.CurrencyMaster{
		{ID: JPYID, Code: money.JPY},
		{ID: USDID, Code: money.USD},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
