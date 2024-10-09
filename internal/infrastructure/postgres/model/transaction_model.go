package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TransactionModel struct {
	bun.BaseModel     `bun:"table:transactions"`
	ID                string    `bun:"id,pk,type:char(26),notnull"`
	AccountID         string    `bun:"account_id,notnull"`
	ReceiverAccountID *string   `bun:"receiver_account_id"`
	Type              string    `bun:"type,type:varchar(20),notnull"`
	Amount            float64   `bun:"amount,type:float8,notnull"`
	CurrencyID        string    `bun:"currency_id,notnull"`
	TransactionAt     time.Time `bun:"transaction_at,notnull"`

	SenderAccount   *AccountModel               `bun:"rel:belongs-to,join:account_id=id"`
	ReceiverAccount *AccountModel               `bun:"rel:belongs-to,join:receiver_account_id=id"`
	Currency        *CurrencyMasterModel        `bun:"rel:belongs-to,join:currency_id=id"`
	TransactionType *TransactionTypeMasterModel `bun:"rel:belongs-to,join:type=type"`
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

var TransactionTypeFK = ForeignKey{
	Table:            "transactions",
	ConstraintName:   "fk_transaction_type",
	Column:           "type",
	ReferencedTable:  "transaction_type_master",
	ReferencedColumn: "type",
}

var TransactionSenderAccountIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*TransactionModel)(nil)).
			Index("transaction_account_id_idx").
			Column("account_id")
	},
}

var TransactionReceiverAccountIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*TransactionModel)(nil)).
			Index("transaction_receiver_account_id_idx").
			Column("receiver_account_id")
	},
}
