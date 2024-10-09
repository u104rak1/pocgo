package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

const (
	jpyID = "01J9R7YPV1FH1V0PPKVSB5C9TQ"
)

func SaveCurrencyMaster(db *bun.DB) error {
	data := []model.CurrencyMasterModel{
		{ID: jpyID, Code: "JPY"},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
