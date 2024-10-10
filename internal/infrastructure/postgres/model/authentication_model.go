package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Authentication struct {
	bun.BaseModel `bun:"table:authentications"`
	UserID        string    `bun:"user_id,pk,notnull"`
	PasswordHash  string    `bun:"password_hash,notnull"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

var AuthenticationUserFK = ForeignKey{
	Table:            "authentications",
	ConstraintName:   "fk_auth_user_id",
	Column:           "user_id",
	ReferencedTable:  "users",
	ReferencedColumn: "id",
}
