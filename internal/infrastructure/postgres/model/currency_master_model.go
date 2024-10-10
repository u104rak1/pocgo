package model

import (
	"github.com/uptrace/bun"
)

type CurrencyMasterModel struct {
	bun.BaseModel `bun:"table:currency_master"`
	ID            string `bun:"id,pk,type:char(26),notnull"`
	Code          string `bun:"code,type:varchar(3),notnull"`
}
