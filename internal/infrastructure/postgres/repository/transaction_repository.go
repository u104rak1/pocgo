package repository

import (
	"context"
	"fmt"

	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type transactionRepository struct {
	*Repository[model.Transaction]
}

func NewTransactionRepository(db *bun.DB) transactionDomain.ITransactionRepository {
	return &transactionRepository{Repository: NewRepository[model.Transaction](db)}
}

func (r *transactionRepository) Save(ctx context.Context, transaction *transactionDomain.Transaction) error {
	currencyCode := transaction.TransferAmount().Currency()

	var currencyID string
	err := r.execDB(ctx).NewSelect().
		Model((*model.CurrencyMaster)(nil)).
		Column("id").
		Where("code = ?", currencyCode).
		Scan(ctx, &currencyID)

	if err != nil {
		return err
	}

	transactionModel := &model.Transaction{
		ID:                transaction.ID(),
		AccountID:         transaction.AccountID(),
		ReceiverAccountID: transaction.ReceiverAccountID(),
		OperationType:     transaction.OperationType(),
		Amount:            transaction.TransferAmount().Amount(),
		CurrencyID:        currencyID,
		TransactionAt:     transaction.TransactionAt(),
	}
	_, err = r.execDB(ctx).NewInsert().Model(transactionModel).Exec(ctx)
	return err
}

func (r *transactionRepository) ListWithTotalByAccountID(ctx context.Context, params transactionDomain.ListTransactionsParams) (transactions []*transactionDomain.Transaction, total int, err error) {
	totalCountQuery := r.execDB(ctx).NewSelect().Model(&model.Transaction{})
	r.buildListQuery(totalCountQuery, params)

	total, err = totalCountQuery.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total transactions: %w", err)
	}

	var transactionModels = []model.Transaction{}
	getQuery := r.execDB(ctx).NewSelect().Model(&transactionModels)
	r.buildListQuery(getQuery, params)

	if *params.Sort == "ASC" {
		getQuery.Order("transaction_at ASC")
	} else {
		getQuery.Order("transaction_at DESC")
	}
	if params.Limit != nil {
		getQuery.Limit(*params.Limit)
	}
	if params.Page != nil && params.Limit != nil {
		getQuery.Offset((*params.Page - 1) * *params.Limit)
	}

	if err := getQuery.Scan(ctx); err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve transactions: %w", err)
	}

	transactions = make([]*transactionDomain.Transaction, len(transactionModels))
	for i, m := range transactionModels {
		transaction, err := transactionDomain.New(
			m.ID,
			m.AccountID,
			m.ReceiverAccountID,
			m.OperationType,
			m.Amount,
			m.Currency.Code,
			m.TransactionAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions[i] = transaction
	}

	return transactions, total, nil
}

func (r *transactionRepository) buildListQuery(query *bun.SelectQuery, params transactionDomain.ListTransactionsParams) {
	query.Relation("Currency").Where("account_id = ?", params.AccountID)

	if params.From != nil {
		query.Where("transaction_at >= ?", *params.From)
	}

	if params.To != nil {
		query.Where("transaction_at <= ?", *params.To)
	}

	if len(params.OperationTypes) > 0 {
		query.Where("operation_type IN (?)", bun.In(params.OperationTypes))
	}
}
