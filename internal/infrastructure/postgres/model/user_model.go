package model

import "github.com/uptrace/bun"

type UserModel struct {
	bun.BaseModel `bun:"table:users"`
	ID            string `bun:"id,pk,notnull"`
	Name          string `bun:"name,type:varchar(20),notnull"`
	Email         string `bun:"email,notnull"`
}

var UserEmailIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*UserModel)(nil)).
			Index("user_email_idx").
			Unique().
			Column("email")
	},
}