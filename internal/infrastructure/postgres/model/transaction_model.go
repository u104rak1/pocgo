package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel     `bun:"table:transactions"`
	ID                string    `bun:"id,pk,type:char(26),notnull"`
	AccountID         string    `bun:"account_id,type:char(26),notnull"`
	ReceiverAccountID *string   `bun:"receiver_account_id,type:char(26)"`
	OperationType     string    `bun:"operation_type,type:varchar(20),notnull"`
	Amount            float64   `bun:"amount,type:float8,notnull"`
	CurrencyID        string    `bun:"currency_id,type:char(26),notnull"`
	TransactionAt     time.Time `bun:"transaction_at,notnull"`

	SenderAccount       *Account             `bun:"rel:belongs-to,join:account_id=id"`
	ReceiverAccount     *Account             `bun:"rel:belongs-to,join:receiver_account_id=id"`
	Currency            *CurrencyMaster      `bun:"rel:belongs-to,join:currency_id=id"`
	OperationTypeMaster *OperationTypeMaster `bun:"rel:belongs-to,join:operation_type=type"`
}

var TransactionAccountFK = ForeignKey{
	Table:            "transactions",
	ConstraintName:   "fk_transaction_account_id",
	Column:           "account_id",
	ReferencedTable:  "accounts",
	ReferencedColumn: "id",
}

var TransactionReceiverAccountFK = ForeignKey{
	Table:            "transactions",
	ConstraintName:   "fk_transaction_receiver_account_id",
	Column:           "receiver_account_id",
	ReferencedTable:  "accounts",
	ReferencedColumn: "id",
}

var TransactionCurrencyFK = ForeignKey{
	Table:            "transactions",
	ConstraintName:   "fk_transaction_currency_id",
	Column:           "currency_id",
	ReferencedTable:  "currency_master",
	ReferencedColumn: "id",
}

var OperationTypeFK = ForeignKey{
	Table:            "transactions",
	ConstraintName:   "fk_transaction_operation_type",
	Column:           "operation_type",
	ReferencedTable:  "operation_type_master",
	ReferencedColumn: "type",
}

var TransactionSenderAccountIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Transaction)(nil)).
			Index("transaction_account_id_idx").
			Column("account_id")
	},
}

var TransactionReceiverAccountIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*Transaction)(nil)).
			Index("transaction_receiver_account_id_idx").
			Column("receiver_account_id")
	},
}
