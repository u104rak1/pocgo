package model

import "github.com/uptrace/bun"

var Models = []interface{}{
	(*CurrencyMasterModel)(nil),
	(*TransactionTypeMasterModel)(nil),
	(*UserModel)(nil),
	(*AccountModel)(nil),
	(*TransactionModel)(nil),
	(*AuthenticationModel)(nil),
}

type IndexQueryCreators func(db *bun.DB) *bun.CreateIndexQuery

func AllIdxCreators() []IndexQueryCreators {
	return append(
		append(AccountUserIDIdxCreator, UserEmailIdxCreator...),
		append(TransactionSenderAccountIDIdxCreator, TransactionReceiverAccountIDIdxCreator...)...,
	)
}

type ForeignKey struct {
	Table            string
	ConstraintName   string
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

var ForeignKeys = []ForeignKey{
	AccountUserFK,
	AccountCurrencyFK,
	AuthenticationUserFK,
	TransactionAccountFK,
	TransactionReceiverAccountFK,
	TransactionCurrencyFK,
	TransactionTypeFK,
}
