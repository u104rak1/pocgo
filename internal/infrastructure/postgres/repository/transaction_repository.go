package repository

import (
	"context"

	transactionDomain "github.com/ucho456job/pocgo/internal/domain/transaction"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
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

func (r *transactionRepository) ListByAccountID(ctx context.Context, accountID string, limit, offset *int) ([]*transactionDomain.Transaction, error) {
	var transactionModels []model.Transaction
	query := r.execDB(ctx).NewSelect().
		Model(&transactionModels).
		Relation("Currency").
		Where("account_id = ?", accountID).
		Order("transaction_at DESC")

	if limit != nil {
		query.Limit(*limit)
	}
	if offset != nil {
		query.Offset(*offset)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	transactions := make([]*transactionDomain.Transaction, len(transactionModels))
	for i, model := range transactionModels {
		transaction, err := transactionDomain.New(
			model.ID,
			model.AccountID,
			model.ReceiverAccountID,
			model.OperationType,
			model.Amount,
			model.Currency.Code,
			model.TransactionAt,
		)
		if err != nil {
			return nil, err
		}
		transactions[i] = transaction
	}

	return transactions, nil
}
