package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TransactionModel struct {
	bun.BaseModel     `bun:"table:transactions"`
	ID                string    `bun:"id,pk,notnull"`
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
