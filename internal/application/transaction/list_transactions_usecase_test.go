package transaction_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	transactionUC "github.com/u104rak1/pocgo/internal/application/transaction"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestListTransactionsUsecase(t *testing.T) {
	type Mocks struct {
		accountServ     *domainMock.MockIAccountService
		transactionServ *domainMock.MockITransactionService
	}

	var (
		userID      = idVO.NewUserIDForTest("user")
		accountID   = idVO.NewAccountIDForTest("account")
		accountName = "test"
		password    = "1234"
		amount      = 1000.0
		currency    = moneyVO.JPY
		time        = timer.GetFixedDate()
		arg         = gomock.Any()
		sort        = "DESC"
		limit       = 10
		page        = 1
	)
	passwordHash, err := passwordUtil.Encode(password)
	assert.NoError(t, err)

	happyCmd := transactionUC.ListTransactionsCommand{
		UserID:         userID.String(),
		AccountID:      accountID.String(),
		From:           &time,
		To:             &time,
		OperationTypes: []string{transactionDomain.Deposit, transactionDomain.Withdraw},
		Sort:           &sort,
		Limit:          &limit,
		Page:           &page,
	}

	tests := []struct {
		caseName string
		cmd      transactionUC.ListTransactionsCommand
		prepare  func(mocks Mocks, account *accountDomain.Account)
		wantErr  bool
	}{
		{
			caseName: "Positive: 取引履歴の取得が成功する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, nil).Return(account, nil)

				tx1, err := transactionDomain.New(account.ID(), nil, transactionDomain.Deposit, amount, currency, time)
				assert.NoError(t, err)

				transactions := []transactionDomain.Transaction{*tx1}
				total := len(transactions)

				mocks.transactionServ.EXPECT().ListWithTotal(arg, arg).Return(transactions, total, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: ユーザーIDが不正な形式である",
			cmd: transactionUC.ListTransactionsCommand{
				UserID: "invalid",
			},
			prepare: func(mocks Mocks, account *accountDomain.Account) {},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座IDが不正な形式である",
			cmd: transactionUC.ListTransactionsCommand{
				UserID:    userID.String(),
				AccountID: "invalid",
			},
			prepare: func(mocks Mocks, account *accountDomain.Account) {},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座認証に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, nil).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 取引履歴の取得に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, nil).Return(account, nil)
				mocks.transactionServ.EXPECT().ListWithTotal(arg, arg).Return(nil, 0, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				accountServ:     domainMock.NewMockIAccountService(ctrl),
				transactionServ: domainMock.NewMockITransactionService(ctrl),
			}

			uc := transactionUC.NewListTransactionsUsecase(
				mocks.accountServ, mocks.transactionServ,
			)
			ctx := context.Background()
			acc, err := accountDomain.Reconstruct(
				accountID.String(), userID.String(), accountName, passwordHash, currency, 0.0, time,
			)
			assert.NoError(t, err)
			tt.prepare(mocks, acc)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.NotEmpty(t, dto.Total)
				assert.NotEmpty(t, dto.Transactions)
				for _, tx := range dto.Transactions {
					assert.NotEmpty(t, tx.ID)
					assert.NotEmpty(t, tx.AccountID)
					assert.NotEmpty(t, tx.OperationType)
					assert.NotEmpty(t, tx.Amount)
					assert.NotEmpty(t, tx.Currency)
					assert.NotEmpty(t, tx.TransactionAt)
				}
			}
		})
	}
}
