package model

import "github.com/uptrace/bun"

type UserModel struct {
	bun.BaseModel `bun:"table:users"`
	ID            string `bun:"id,pk,type:char(26),notnull"`
	Name          string `bun:"name,type:varchar(20),notnull"`
	Email         string `bun:"email,notnull"`

	Authentication *AuthenticationModel `bun:"rel:has-one,join:id=user_id"`
	Accounts       []*AccountModel      `bun:"rel:has-many,join:id=user_id"`
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
