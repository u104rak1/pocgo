package repository_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestTransactionRepository_Save(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewTransactionRepository)

	accountID := idVO.NewAccountIDForTest("account")
	money, err := moneyVO.New(1000, moneyVO.JPY)
	assert.NoError(t, err)

	transactionAt := timer.GetFixedDate()
	transaction, err := transactionDomain.New(accountID, nil, transactionDomain.Deposit, money.Amount(), money.Currency(), transactionAt)
	assert.NoError(t, err)

	currencyID := idVO.GenerateStaticULID(moneyVO.JPY)
	currencySelectQuery := `SELECT "currency_master"."id" FROM "currency_master" WHERE (code = 'JPY')`

	expectQuery := fmt.Sprintf(
		`INSERT INTO "transactions" ("id", "account_id", "receiver_account_id", "operation_type", "amount", "currency_id", "transaction_at")
		VALUES ('%s', '%s', DEFAULT, '%s', %.0f, '%s', '%s')
		RETURNING "receiver_account_id"`,
		transaction.IDString(), transaction.AccountIDString(), transaction.OperationType(),
		transaction.TransferAmount().Amount(), currencyID, transactionAt.Format("2006-01-02 15:04:05-07:00"),
	)

	tests := []struct {
		caseName string
		prepare  func()
		wantErr  bool
	}{
		{
			caseName: "Positive: 取引の保存が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(currencyID))
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"receiver_account_id"}))
			},
			wantErr: false,
		},
		{
			caseName: "Negative: 通貨の取得に失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 取引の保存に失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(currencySelectQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(currencyID))
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := repo.Save(ctx, transaction)

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

// func TestTransactionRepository_ListWithTotalByAccountID(t *testing.T) {
// 	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewTransactionRepository)

// 	accountID := idVO.NewAccountIDForTest("account")
// 	now := timer.GetFixedDate()
// 	params := transactionDomain.ListTransactionsParams{
// 		AccountID:      accountID,
// 		From:           &now,
// 		To:             &now,
// 		OperationTypes: []string{transactionDomain.Deposit, transactionDomain.Withdrawal},
// 		Sort:           new(string),
// 		Limit:          new(int),
// 		Page:           new(int),
// 	}

// 	// expectCountQuery := fmt.Sprintf(`
// 	// 	SELECT count(*) FROM "transactions" AS "transaction"
// 	// 	LEFT JOIN "currency_master" AS "currency" ON ("currency"."id" = "transaction"."currency_id")
// 	// 	WHERE (account_id = '%s')
// 	// 	AND (transaction_at >= '%s')
// 	// 	AND (transaction_at <= '%s')
// 	// 	AND (operation_type IN ('DEPOSIT','WITHDRAWAL'))
// 	// `, accountID.String(), params.From.Format("2006-01-02 15:04:05-07:00"),
// 	// 	params.To.Format("2006-01-02 15:04:05-07:00"))

// 	expectSelectQuery := fmt.Sprintf(`
// 		SELECT "id", "account_id", "amount", "operation_type", "description", "transaction_at"
// 		FROM "transactions" AS "transaction"
// 		LEFT JOIN "currency_master" AS "currency" ON ("currency"."id" = "transaction"."currency_id")
// 		WHERE (account_id = '%s')
// 		AND (transaction_at >= '%s')
// 		AND (transaction_at <= '%s')
// 		AND (operation_type IN ('DEPOSIT','WITHDRAWAL'))
// 		ORDER BY transaction_at DESC
// 	`, accountID.String(), params.From.Format("2006-01-02 15:04:05-07:00"),
// 		params.To.Format("2006-01-02 15:04:05-07:00"))

// 	tests := []struct {
// 		caseName    string
// 		params      transactionDomain.ListTransactionsParams
// 		prepare     func()
// 		wantTotal   int
// 		wantTxCount int
// 		wantErr     bool
// 	}{
// 		{
// 			caseName: "Positive: 全てのパラメータが設定された場合",
// 			params:   params,
// 			prepare: func() {
// 				expectCountQuery := fmt.Sprintf(`SELECT count(*) FROM "transactions" AS "transaction" LEFT JOIN "currency_master" AS "currency" ON ("currency"."id" = "transaction"."currency_id") WHERE (account_id = '%s') AND (transaction_at >= '%s') AND (transaction_at <= '%s') AND (operation_type IN ('DEPOSIT','WITHDRAWAL'))`, accountID.String(), params.From.Format("2006-01-02 15:04:05-07:00"), params.To.Format("2006-01-02 15:04:05-07:00"))

// 				mock.ExpectQuery(regexp.QuoteMeta(expectCountQuery)).
// 					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

// 				fmt.Println("Executing count query:", expectCountQuery)

// 				rows := sqlmock.NewRows([]string{
// 					"id", "account_id", "amount", "operation_type", "description", "transaction_at",
// 				})
// 				for i := 0; i < 2; i++ {
// 					rows.AddRow(
// 						fmt.Sprintf("tx_%d", i),
// 						accountID.String(),
// 						1000,
// 						"DEPOSIT",
// 						"Test Transaction",
// 						now,
// 					)
// 				}
// 				mock.ExpectQuery(regexp.QuoteMeta(expectSelectQuery)).
// 					WillReturnRows(rows)
// 			},
// 			wantTotal:   2,
// 			wantTxCount: 2,
// 			wantErr:     false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.caseName, func(t *testing.T) {
// 			tt.prepare()
// 			_, _, err := repo.ListWithTotalByAccountID(ctx, tt.params)

// 			// if tt.wantErr {
// 			// 	assert.Error(t, err)
// 			// 	assert.Equal(t, 0, total)
// 			// 	assert.Len(t, transactions, 0)
// 			// } else {
// 			// 	assert.NoError(t, err)
// 			// 	assert.Equal(t, tt.wantTotal, total)
// 			// 	assert.Len(t, transactions, tt.wantTxCount)
// 			// }
// 			err = mock.ExpectationsWereMet()
// 			assert.NoError(t, err)
// 		})
// 	}
// }
