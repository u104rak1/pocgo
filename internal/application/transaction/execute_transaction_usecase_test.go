package transaction_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	transactionUC "github.com/u104rak1/pocgo/internal/application/transaction"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestExecuteTransactionUsecase(t *testing.T) {
	type Mocks struct {
		accountServ     *domainMock.MockIAccountService
		transactionServ *domainMock.MockITransactionService
	}

	var (
		userID              = idVO.NewUserIDForTest("user")
		accountID           = idVO.NewAccountIDForTest("account")
		receiverID          = idVO.NewAccountIDForTest("receiver")
		accountName         = "test"
		receiverAccountName = "receiver"
		password            = "1234"
		amount              = 1000.0
		currency            = moneyVO.JPY
		time                = timer.GetFixedDate()
		arg                 = gomock.Any()
	)
	passwordHash, err := passwordUtil.Encode(password)
	assert.NoError(t, err)

	happyDepositCmd := transactionUC.ExecuteTransactionCommand{
		UserID:        userID.String(),
		AccountID:     accountID.String(),
		Password:      password,
		OperationType: transactionDomain.Deposit,
		Amount:        amount,
		Currency:      currency,
	}

	happyWithdrawCmd := transactionUC.ExecuteTransactionCommand{
		UserID:        userID.String(),
		AccountID:     accountID.String(),
		Password:      password,
		OperationType: transactionDomain.Withdrawal,
		Amount:        amount,
		Currency:      currency,
	}

	receiverIDStr := receiverID.String()
	happyTransferCmd := transactionUC.ExecuteTransactionCommand{
		UserID:            userID.String(),
		AccountID:         accountID.String(),
		Password:          password,
		OperationType:     transactionDomain.Transfer,
		Amount:            amount,
		Currency:          currency,
		ReceiverAccountID: &receiverIDStr,
	}

	tests := []struct {
		caseName string
		cmd      transactionUC.ExecuteTransactionCommand
		prepare  func(mocks Mocks, account *accountDomain.Account)
		wantErr  bool
	}{
		{
			caseName: "Positive: 入金取引が成功する",
			cmd:      happyDepositCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)

				tx, err := transactionDomain.New(account.ID(), nil, transactionDomain.Deposit, amount, currency, time)
				assert.NoError(t, err)
				mocks.transactionServ.EXPECT().Deposit(arg, arg, arg, arg).Return(tx, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Positive: 出金取引が成功する",
			cmd:      happyWithdrawCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)

				tx, err := transactionDomain.New(account.ID(), nil, transactionDomain.Withdrawal, amount, currency, time)
				assert.NoError(t, err)
				mocks.transactionServ.EXPECT().Withdraw(arg, arg, arg, arg).Return(tx, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Positive: 送金取引が成功する",
			cmd:      happyTransferCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)

				receiverAccount, err := accountDomain.Reconstruct(
					receiverID.String(), userID.String(), receiverAccountName, passwordHash, currency, 0.0, time,
				)
				assert.NoError(t, err)
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, nil, nil).Return(receiverAccount, nil)

				receiverAccountID := receiverAccount.ID()
				tx, err := transactionDomain.New(account.ID(), &receiverAccountID, transactionDomain.Transfer, amount, currency, time)
				assert.NoError(t, err)
				mocks.transactionServ.EXPECT().Transfer(arg, arg, arg, arg, arg).Return(tx, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: ユーザーIDが不正な形式である",
			cmd: transactionUC.ExecuteTransactionCommand{
				UserID: "invalid",
			},
			prepare: func(mocks Mocks, account *accountDomain.Account) {},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座IDが不正な形式である",
			cmd: transactionUC.ExecuteTransactionCommand{
				UserID:    userID.String(),
				AccountID: "invalid",
			},
			prepare: func(mocks Mocks, account *accountDomain.Account) {},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座認証に失敗する",
			cmd:      happyDepositCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 入金処理に失敗する",
			cmd:      happyDepositCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)
				mocks.transactionServ.EXPECT().Deposit(arg, arg, arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 出金処理に失敗する",
			cmd:      happyWithdrawCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)
				mocks.transactionServ.EXPECT().Withdraw(arg, arg, arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 受け取り口座の取得に失敗する",
			cmd:      happyTransferCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, nil, nil).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 送金処理に失敗する",
			cmd:      happyTransferCmd,
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)

				receiverAccount, err := accountDomain.Reconstruct(
					receiverID.String(), userID.String(), receiverAccountName, passwordHash, currency, 0.0, time,
				)
				assert.NoError(t, err)
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, nil, nil).Return(receiverAccount, nil)

				mocks.transactionServ.EXPECT().Transfer(arg, arg, arg, arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: サポートされていない取引種別である",
			cmd: transactionUC.ExecuteTransactionCommand{
				UserID:        userID.String(),
				AccountID:     accountID.String(),
				Password:      password,
				OperationType: "UNSUPPORTED",
				Amount:        amount,
				Currency:      currency,
			},
			prepare: func(mocks Mocks, account *accountDomain.Account) {
				mocks.accountServ.EXPECT().GetAndAuthorize(arg, arg, arg, arg).Return(account, nil)
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
			mockUnitOfWork := &appMock.MockIUnitOfWorkWithResult[transactionDomain.Transaction]{}

			uc := transactionUC.NewExecuteTransactionUsecase(
				mocks.accountServ, mocks.transactionServ, mockUnitOfWork,
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
				assert.NotEmpty(t, dto.ID)
				assert.Equal(t, tt.cmd.AccountID, dto.AccountID)
				assert.Equal(t, tt.cmd.ReceiverAccountID, dto.ReceiverAccountID)
				assert.Equal(t, tt.cmd.OperationType, dto.OperationType)
				assert.Equal(t, tt.cmd.Amount, dto.Amount)
				assert.Equal(t, tt.cmd.Currency, dto.Currency)
				assert.NotEmpty(t, dto.TransactionAt)
			}
		})
	}
}
