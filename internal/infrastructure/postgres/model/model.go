package model

import "github.com/uptrace/bun"

type IndexQueryCreators func(db *bun.DB) *bun.CreateIndexQuery

var Models = []interface{}{
	(*UserModel)(nil),
	(*AccountModel)(nil),
}

func AllIdxCreators() []IndexQueryCreators {
	return append(
		AccountUserIDIdxCreator,
		UserEmailIdxCreator...,
	)
}
