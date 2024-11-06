package model

import (
	"github.com/uptrace/bun"
)

type OperationTypeMaster struct {
	bun.BaseModel `bun:"table:operation_type_master"`
	Type          string `bun:"type,pk,type:varchar(20),notnull"`
}
