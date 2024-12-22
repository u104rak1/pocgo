package repository

import (
	"context"
	"database/sql"
	"errors"

	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type accountRepository struct {
	*Repository[model.Account]
}

func NewAccountRepository(db *bun.DB) accountDomain.IAccountRepository {
	return &accountRepository{Repository: NewRepository[model.Account](db)}
}

func (r *accountRepository) Save(ctx context.Context, account *accountDomain.Account) error {
	currencyCode := account.Balance().Currency()

	var currencyID string
	err := r.ExecDB(ctx).NewSelect().
		Model((*model.CurrencyMaster)(nil)).
		Column("id").
		Where("code = ?", currencyCode).
		Scan(ctx, &currencyID)

	if err != nil {
		return err
	}

	accountModel := &model.Account{
		ID:           account.IDString(),
		Name:         account.Name(),
		UserID:       account.UserIDString(),
		PasswordHash: account.PasswordHash(),
		Balance:      account.Balance().Amount(),
		CurrencyID:   currencyID,
		UpdatedAt:    account.UpdatedAt(),
	}

	// TODO: If use a subquery, the following error will occur, so first get the current_id and then update it.
	// pgdriver.Error: ERROR: insert or update on table "accounts" violates foreign key constraint "fk_account_currency_id" (SQLSTATE=23503)
	_, err = r.ExecDB(ctx).NewInsert().Model(accountModel).On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("user_id = EXCLUDED.user_id").
		Set("password_hash = EXCLUDED.password_hash").
		Set("balance = EXCLUDED.balance").
		Set("currency_id = EXCLUDED.currency_id").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

func (r *accountRepository) FindByID(ctx context.Context, id idVO.AccountID) (*accountDomain.Account, error) {
	accountModel := &model.Account{}

	if err := r.ExecDB(ctx).NewSelect().
		Model(accountModel).
		Relation("Currency").
		Where("account.id = ?", id.String()).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return accountDomain.Reconstruct(
		accountModel.ID,
		accountModel.UserID,
		accountModel.Name,
		accountModel.PasswordHash,
		accountModel.Currency.Code,
		accountModel.Balance,
		accountModel.UpdatedAt,
	)
}

func (r *accountRepository) CountByUserID(ctx context.Context, userID idVO.UserID) (int, error) {
	return r.ExecDB(ctx).NewSelect().Model((*model.Account)(nil)).Where("user_id = ?", userID.String()).Count(ctx)
}
