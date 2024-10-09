package model

import "github.com/uptrace/bun"

type AuthenticationModel struct {
	bun.BaseModel `bun:"table:authentications"`
	ID            string `bun:"id,pk,notnull"`
	UserID        string `bun:"user_id,unique,notnull"`
	PasswordHash  string `bun:"password_hash,notnull"`
}

var AuthenticationUserFK = ForeignKey{
	Table:            "authentications",
	ConstraintName:   "fk_auth_user_id",
	Column:           "user_id",
	ReferencedTable:  "users",
	ReferencedColumn: "id",
}

var AuthenticationUserIDIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*AuthenticationModel)(nil)).
			Index("authentication_user_id_idx").
			Unique().
			Column("user_id")
	},
}
