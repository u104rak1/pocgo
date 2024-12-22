package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
)

func TestAccountRepository_Save(t *testing.T) {
	repo, mock, ctx, db := PrepareTestRepository(t, repository.NewAccountRepository)
	userID := idVO.NewUserIDForTest("user")
	money, err := moneyVO.New(1000, moneyVO.JPY)
	assert.NoError(t, err)
	account, err := accountDomain.New(userID, money.Amount(), "Test Account", "1234", money.Currency())
	assert.NoError(t, err)
	currencyID := idVO.GenerateStaticULID("JPY")

	currencySelectQuery := `SELECT "currency_master"."id" FROM "currency_master" WHERE (code = 'JPY')`
	expectInsertQuery := fmt.Sprintf(`
		INSERT INTO "accounts" AS "account" ("id", "user_id", "name", "password_hash", "balance", "currency_id", "updated_at", "deleted_at")
		VALUES ('%s', '%s', '%s', '%s', %.0f, '%s', '%s', DEFAULT)
		ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		user_id = EXCLUDED.user_id,
		password_hash = EXCLUDED.password_hash,
		balance = EXCLUDED.balance,
		currency_id = EXCLUDED.currency_id,
		updated_at = EXCLUDED.updated_at
		RETURNING "deleted_at"
	`, account.IDString(), account.UserIDString(), account.Name(), account.PasswordHash(), account.Balance().Amount(), currencyID, account.UpdatedAt().Format("2006-01-02 15:04:05-07:00"))

	tests := []struct {
		caseName string
		prepare  func()
		wantErr  bool
	}{
		{
			caseName: "Positive: アカウントの保存が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(currencyID))
				mock.ExpectQuery(regexp.QuoteMeta(expectInsertQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
			},
			wantErr: false,
		},
		{
			caseName: "Negative: 通貨マスタの取得に失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: アカウントの保存に失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(currencyID))
				mock.ExpectQuery(regexp.QuoteMeta(expectInsertQuery)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := repo.Save(ctx, account)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}

	unitOfWork := repository.NewUnitOfWork(db)
	testsWithTx := []struct {
		caseName string
		prepare  func()
		wantErr  bool
	}{
		{
			caseName: "Positive: トランザクション内でアカウントの保存が成功する",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(currencyID))
				mock.ExpectQuery(regexp.QuoteMeta(expectInsertQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			caseName: "Negative: SQLエラーでロールバックされる",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range testsWithTx {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
				return repo.Save(ctx, account)
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAccountRepository_FindByID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewAccountRepository)
	userID := idVO.NewUserIDForTest("user_id_1")
	money, err := moneyVO.New(1000, moneyVO.JPY)
	assert.NoError(t, err)
	account, err := accountDomain.New(userID, money.Amount(), "Test Account", "1234", money.Currency())
	assert.NoError(t, err)
	currencyID := idVO.GenerateStaticULID("JPY")

	expectQuery := fmt.Sprintf(`
		SELECT "account"."id", "account"."user_id", "account"."name", "account"."password_hash",
		"account"."balance", "account"."currency_id", "account"."updated_at", "account"."deleted_at",
		"currency"."id" AS "currency__id", "currency"."code" AS "currency__code"
		FROM "accounts" AS "account"
		LEFT JOIN "currency_master" AS "currency" ON ("currency"."id" = "account"."currency_id")
		WHERE (account.id = '%s') AND "account"."deleted_at" IS NULL
	`, account.IDString())

	tests := []struct {
		caseName    string
		prepare     func()
		wantAccount *accountDomain.Account
		wantErr     bool
	}{
		{
			caseName: "Positive: IDでアカウント取得が成功する",
			prepare: func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "user_id", "password_hash", "balance", "currency_id",
					"updated_at", "deleted_at", "currency__id", "currency__code",
				}).AddRow(
					account.IDString(), account.Name(), account.UserIDString(),
					account.PasswordHash(), account.Balance().Amount(), currencyID,
					account.UpdatedAt(), nil, currencyID, "JPY",
				)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantAccount: account,
			wantErr:     false,
		},
		{
			caseName: "Positive: アカウントが見つからない場合、nilを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantAccount: nil,
			wantErr:     false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantAccount: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundAccount, err := repo.FindByID(ctx, account.ID())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, foundAccount)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAccount, foundAccount)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAccountRepository_CountByUserID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewAccountRepository)
	userID := idVO.NewUserIDForTest("user_id_1")

	expectQuery := fmt.Sprintf(`
		SELECT count(*) FROM "accounts" AS "account"
		WHERE (user_id = '%s') AND "account"."deleted_at" IS NULL
	`, userID.String())

	tests := []struct {
		caseName  string
		prepare   func()
		wantCount int
		wantErr   bool
	}{
		{
			caseName: "Positive: ユーザーIDでアカウント数の取得が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).
					WillReturnError(assert.AnError)
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			count, err := repo.CountByUserID(ctx, userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
