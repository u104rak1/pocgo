package model

import (
	"time"

	"github.com/uptrace/bun"
)

type AccountModel struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            string    `bun:"id,pk,notnull"`
	UserID        string    `bun:"user_id,notnull"`
	Name          string    `bun:"name,type:varchar(10),notnull"`
	Password      string    `bun:"password,notnull"`
	Balance       float64   `bun:"balance,type:float8,notnull"`
	Currency      string    `bun:"currency,type:varchar(3),notnull"`
	LastUpdatedAt time.Time `bun:"last_updated_at"`
}

var AccountUserIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*AccountModel)(nil)).
			Index("account_user_id_idx").
			Unique().
			Column("user_id")
	},
}
