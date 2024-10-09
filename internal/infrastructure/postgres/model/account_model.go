package model

import (
	"time"

	"github.com/uptrace/bun"
)

type AccountModel struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            string    `bun:"id,pk,type:char(26),notnull"`
	UserID        string    `bun:"user_id,notnull"`
	Name          string    `bun:"name,type:varchar(10)"`
	PasswordHash  string    `bun:"password_hash,notnull"`
	Balance       float64   `bun:"balance,type:float8,notnull"`
	CurrencyID    string    `bun:"currency_id,notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at,notnull"`

	User                 *UserModel           `bun:"rel:belongs-to,join:user_id=id"`
	SentTransactions     []*TransactionModel  `bun:"rel:has-many,join:id=account_id"`
	ReceivedTransactions []*TransactionModel  `bun:"rel:has-many,join:id=receiver_account_id"`
	Currency             *CurrencyMasterModel `bun:"rel:belongs-to,join:currency_id=id"`
}

var AccountUserFK = ForeignKey{
	Table:            "accounts",
	ConstraintName:   "fk_account_user_id",
	Column:           "user_id",
	ReferencedTable:  "users",
	ReferencedColumn: "id",
}

var AccountCurrencyFK = ForeignKey{
	Table:            "accounts",
	ConstraintName:   "fk_account_currency_id",
	Column:           "currency_id",
	ReferencedTable:  "currency_master",
	ReferencedColumn: "id",
}

var AccountUserIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*AccountModel)(nil)).
			Index("account_user_id_idx").
			Column("user_id")
	},
}
