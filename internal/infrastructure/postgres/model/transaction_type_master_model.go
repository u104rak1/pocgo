package model

import (
	"github.com/uptrace/bun"
)

type TransactionTypeMaster struct {
	bun.BaseModel `bun:"table:transaction_type_master"`
	Type          string `bun:"type,pk,type:varchar(20),notnull"`
}
