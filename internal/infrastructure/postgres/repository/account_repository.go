package repository

import (
	"context"
	"time"

	accountDomain "github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

type IAccountRepository interface {
	Save(ctx context.Context, account *accountDomain.Account) error
	FindByID(ctx context.Context, id string) (*accountDomain.Account, error)
	ListByUserID(ctx context.Context, userID string) ([]*accountDomain.Account, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, id string) error
}

type accountRepository struct {
	db *bun.DB
}

func NewAccountRepository(db *bun.DB) IAccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Save(ctx context.Context, account *accountDomain.Account) error {
	accountModel := &model.Account{
		ID:           account.ID(),
		Name:         account.Name(),
		UserID:       account.UserID(),
		PasswordHash: account.PasswordHash(),
		Balance:      account.Balance().Amount(),
		UpdatedAt:    account.UpdatedAt(),
	}
	currencyCode := account.Balance().Currency()

	_, err := r.db.NewInsert().Model(accountModel).On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("user_id = EXCLUDED.user_id").
		Set("password_hash = EXCLUDED.password_hash").
		Set("balance = EXCLUDED.balance").
		Set("currency_id = (SELECT id FROM currency_master WHERE code = ?)", currencyCode).
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	return err
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (*accountDomain.Account, error) {
	accountModel := &model.Account{}
	var currencyCode string

	if err := r.db.NewSelect().
		Model(accountModel).
		ColumnExpr("account.*, currency_master.code AS currency_code").
		Join("JOIN currency_master ON currency_master.id = account.currency_id").
		Where("account.id = ?", id).
		Scan(ctx, &currencyCode); err != nil {
		return nil, err
	}

	return accountDomain.Reconstruct(
		accountModel.ID,
		accountModel.UserID,
		accountModel.Name,
		accountModel.PasswordHash,
		accountModel.Balance,
		currencyCode,
		accountModel.UpdatedAt,
	)
}

func (r *accountRepository) ListByUserID(ctx context.Context, userID string) ([]*accountDomain.Account, error) {
	var accountModels []*model.Account
	var currencyCodes []string

	if err := r.db.NewSelect().
		Model(&accountModels).
		ColumnExpr("account.*, currency_master.code AS currency_code").
		Join("JOIN currency_master ON currency_master.id = account.currency_id").
		Where("account.user_id = ?", userID).
		Scan(ctx, &currencyCodes); err != nil {
		return nil, err
	}

	accounts := make([]*accountDomain.Account, len(accountModels))

	for i, accountModel := range accountModels {
		account, err := accountDomain.Reconstruct(
			accountModel.ID,
			accountModel.UserID,
			accountModel.Name,
			accountModel.PasswordHash,
			accountModel.Balance,
			currencyCodes[i],
			accountModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts[i] = account
	}

	return accounts, nil
}

func (r *accountRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	return r.db.NewSelect().Model((*model.Account)(nil)).Where("user_id = ?", userID).Count(ctx)
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewUpdate().
		Model(&model.Account{ID: id, DeletedAt: time.Now()}).
		Column("deleted_at").
		WherePK().
		Exec(ctx)
	return err
}
