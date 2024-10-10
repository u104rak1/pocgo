package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            string    `bun:"id,pk,type:char(26),notnull"`
	Name          string    `bun:"name,type:varchar(20),notnull"`
	Email         string    `bun:"email,notnull"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`

	Authentication *Authentication `bun:"rel:has-one,join:id=user_id"`
	Accounts       []*Account      `bun:"rel:has-many,join:id=user_id"`
}

var UserEmailIdxCreator = []IndexQueryCreators{
	func(db *bun.DB) *bun.CreateIndexQuery {
		return db.NewCreateIndex().
			Model((*User)(nil)).
			Index("user_email_idx").
			Unique().
			Column("email")
	},
}
