package model

import "github.com/uptrace/bun"

var Models = []interface{}{
	(*CurrencyMaster)(nil),
	(*OperationTypeMaster)(nil),
	(*User)(nil),
	(*Account)(nil),
	(*Transaction)(nil),
	(*Authentication)(nil),
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
	OperationTypeFK,
}
