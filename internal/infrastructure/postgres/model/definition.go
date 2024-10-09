package model

import "github.com/uptrace/bun"

type IndexQueryCreators func(db *bun.DB) *bun.CreateIndexQuery

var Models = []interface{}{
	(*CurrencyMasterModel)(nil),
	(*TransactionTypeMasterModel)(nil),
	(*UserModel)(nil),
	(*AccountModel)(nil),
	(*TransactionModel)(nil),
	(*AuthenticationModel)(nil),
}

func AllIdxCreators() []IndexQueryCreators {
	return append(
		append(
			append(AccountUserIDIdxCreator, UserEmailIdxCreator...),
			AuthenticationUserIDIdxCreator...,
		),
		append(TransactionSenderAccountIDIdxCreator, TransactionReceiverAccountIDIdxCreator...)...,
	)
}
