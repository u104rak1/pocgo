package transaction_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/numutil"
	"github.com/u104rak1/pocgo/pkg/strutil"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestDeposit(t *testing.T) {
	type Mocks struct {
		accountRepo     *mock.MockIAccountRepository
		transactionRepo *mock.MockITransactionRepository
	}

	var (
		userID        = idVO.NewUserIDForTest("user")
		name          = "account-name"
		password      = "1234"
		balance       = 100.0
		currency      = moneyVO.JPY
		depositAmount = 50.0
		arg           = gomock.Any()
	)

	tests := []struct {
		caseName string
		account  *accountDomain.Account
		amount   float64
		currency string
		setup    func(mocks Mocks)
		errMsg   string
	}{
		{
			caseName: "Positive: 入金が成功する",
			amount:   depositAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: money.Depositが失敗した場合はエラーが返る（通貨単位が異なる）",
			amount:   depositAmount,
			currency: moneyVO.USD,
			setup:    func(mocks Mocks) {},
			errMsg:   moneyVO.ErrAddDifferentCurrency.Error(),
		},
		{
			caseName: "Negative: 口座の保存が失敗した場合はエラーが返る",
			amount:   depositAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: 取引の保存が失敗した場合はエラーが返る",
			amount:   depositAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		// 取引の作成を意図的に失敗させるのが難しいので、テストを省略する
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				accountRepo:     mock.NewMockIAccountRepository(ctrl),
				transactionRepo: mock.NewMockITransactionRepository(ctrl),
			}
			service := transactionDomain.NewService(mocks.accountRepo, mocks.transactionRepo)
			ctx := context.Background()
			tt.setup(mocks)
			account, err := accountDomain.New(userID, balance, name, password, currency)
			assert.NoError(t, err)

			transaction, err := service.Deposit(ctx, account, tt.amount, tt.currency)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Empty(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, account.ID(), transaction.AccountID())
				assert.Equal(t, tt.amount, transaction.TransferAmount().Amount())
				assert.Equal(t, tt.currency, transaction.TransferAmount().Currency())
				assert.Equal(t, "DEPOSIT", transaction.OperationType())
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	type Mocks struct {
		accountRepo     *mock.MockIAccountRepository
		transactionRepo *mock.MockITransactionRepository
	}

	var (
		userID         = idVO.NewUserIDForTest("user")
		name           = "account-name"
		password       = "1234"
		balance        = 100.0
		currency       = moneyVO.JPY
		withdrawAmount = 50.0
		arg            = gomock.Any()
	)

	tests := []struct {
		caseName string
		account  *accountDomain.Account
		amount   float64
		currency string
		setup    func(mocks Mocks)
		errMsg   string
	}{
		{
			caseName: "Positive: 出金が成功する",
			amount:   withdrawAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: money.Withdrawが失敗した場合はエラーが返る（通貨単位が異なる）",
			amount:   withdrawAmount,
			currency: moneyVO.USD,
			setup:    func(mocks Mocks) {},
			errMsg:   moneyVO.ErrSubDifferentCurrency.Error(),
		},
		{
			caseName: "Negative: 口座の保存が失敗した場合はエラーが返る",
			amount:   withdrawAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: 取引の保存が失敗した場合はエラーが返る",
			amount:   withdrawAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				accountRepo:     mock.NewMockIAccountRepository(ctrl),
				transactionRepo: mock.NewMockITransactionRepository(ctrl),
			}
			service := transactionDomain.NewService(mocks.accountRepo, mocks.transactionRepo)
			ctx := context.Background()
			tt.setup(mocks)
			account, err := accountDomain.New(userID, balance, name, password, currency)
			assert.NoError(t, err)

			transaction, err := service.Withdraw(ctx, account, tt.amount, tt.currency)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Empty(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, account.ID(), transaction.AccountID())
				assert.Equal(t, tt.amount, transaction.TransferAmount().Amount())
				assert.Equal(t, tt.currency, transaction.TransferAmount().Currency())
				assert.Equal(t, "WITHDRAW", transaction.OperationType())
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	type Mocks struct {
		accountRepo     *mock.MockIAccountRepository
		transactionRepo *mock.MockITransactionRepository
	}

	var (
		userID         = idVO.NewUserIDForTest("user")
		name           = "account-name"
		password       = "1234"
		balance        = 100.0
		currency       = moneyVO.JPY
		transferAmount = 50.0
		arg            = gomock.Any()
	)

	tests := []struct {
		caseName string
		amount   float64
		currency string
		setup    func(mocks Mocks)
		errMsg   string
	}{
		{
			caseName: "Positive: 送金が成功する",
			amount:   transferAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil).Times(2)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: money.Depositが失敗した場合はエラーが返る（通貨単位が異なる）",
			amount:   transferAmount,
			currency: moneyVO.USD,
			setup:    func(mocks Mocks) {},
			errMsg:   moneyVO.ErrAddDifferentCurrency.Error(),
		},
		{
			caseName: "Negative: money.Withdrawが失敗した場合はエラーが返る（送金元の残高不足）",
			amount:   balance + 1,
			currency: moneyVO.JPY,
			setup:    func(mocks Mocks) {},
			errMsg:   moneyVO.ErrInsufficientBalance.Error(),
		},
		{
			caseName: "Negative: 送金元口座の保存が失敗した場合はエラーが返る",
			amount:   transferAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: 送金先口座の保存が失敗した場合はエラーが返る",
			amount:   transferAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: 取引の保存が失敗した場合はエラーが返る",
			amount:   transferAmount,
			currency: moneyVO.JPY,
			setup: func(mocks Mocks) {
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil).Times(2)
				mocks.transactionRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				accountRepo:     mock.NewMockIAccountRepository(ctrl),
				transactionRepo: mock.NewMockITransactionRepository(ctrl),
			}
			service := transactionDomain.NewService(mocks.accountRepo, mocks.transactionRepo)
			ctx := context.Background()
			tt.setup(mocks)
			senderAccount, err := accountDomain.New(userID, balance, name, password, currency)
			assert.NoError(t, err)
			receiverAccount, err := accountDomain.New(userID, balance, name, password, currency)
			assert.NoError(t, err)

			transaction, err := service.Transfer(ctx, senderAccount, receiverAccount, tt.amount, tt.currency)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Empty(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, senderAccount.ID(), transaction.AccountID())
				assert.Equal(t, receiverAccount.ID(), *transaction.ReceiverAccountID())
				assert.Equal(t, tt.amount, transaction.TransferAmount().Amount())
				assert.Equal(t, tt.currency, transaction.TransferAmount().Currency())
				assert.Equal(t, "TRANSFER", transaction.OperationType())
			}
		})
	}
}

func TestListWithTotal(t *testing.T) {
	type Mocks struct {
		accountRepo     *mock.MockIAccountRepository
		transactionRepo *mock.MockITransactionRepository
	}

	var (
		accountID = idVO.NewAccountIDForTest("account")
		arg       = gomock.Any()
	)

	tests := []struct {
		caseName         string
		params           transactionDomain.ListTransactionsParams
		wantTransactions []*transactionDomain.Transaction
		wantTotal        int
		setup            func(mocks Mocks, transactions []*transactionDomain.Transaction)
		errMsg           string
	}{
		{
			caseName: "Positive: 取引一覧が取得できる（パラメータを全て指定）",
			params: transactionDomain.ListTransactionsParams{
				AccountID: accountID,
				From:      timer.TimePointer(timer.GetFixedDate()),
				To:        timer.TimePointer(timer.GetFixedDate()),
				OperationTypes: []string{
					transactionDomain.Deposit,
					transactionDomain.Withdraw,
					transactionDomain.Transfer,
				},
				Sort:  strutil.StrPointer("ASC"),
				Limit: numutil.IntPointer(20),
				Page:  numutil.IntPointer(2),
			},
			wantTotal: 2,
			setup: func(mocks Mocks, transactions []*transactionDomain.Transaction) {
				mocks.transactionRepo.EXPECT().
					ListWithTotalByAccountID(arg, arg).Return(transactions, 2, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Positive: 取引一覧が取得できる（最低限のパラメータを指定）",
			params: transactionDomain.ListTransactionsParams{
				AccountID: accountID,
			},
			wantTotal: 2,
			setup: func(mocks Mocks, transactions []*transactionDomain.Transaction) {
				mocks.transactionRepo.EXPECT().
					ListWithTotalByAccountID(arg, arg).Return(transactions, 2, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: ListWithTotalByAccountIDがエラーを返した場合はエラーが返される",
			params: transactionDomain.ListTransactionsParams{
				AccountID: accountID,
			},
			wantTransactions: nil,
			wantTotal:        0,
			setup: func(mocks Mocks, transactions []*transactionDomain.Transaction) {
				mocks.transactionRepo.EXPECT().
					ListWithTotalByAccountID(arg, arg).Return(nil, 0, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				accountRepo:     mock.NewMockIAccountRepository(ctrl),
				transactionRepo: mock.NewMockITransactionRepository(ctrl),
			}
			service := transactionDomain.NewService(mocks.accountRepo, mocks.transactionRepo)
			ctx := context.Background()
			tx1, err := transactionDomain.New(
				accountID,
				nil,
				transactionDomain.Deposit,
				1000.0,
				moneyVO.JPY,
				timer.GetFixedDate(),
			)
			assert.NoError(t, err)

			tx2, err := transactionDomain.New(
				accountID,
				nil,
				transactionDomain.Withdraw,
				500.0,
				moneyVO.JPY,
				timer.GetFixedDate(),
			)
			assert.NoError(t, err)

			transactions := []*transactionDomain.Transaction{tx1, tx2}
			tt.setup(mocks, transactions)

			txs, total, err := service.ListWithTotal(ctx, tt.params)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Empty(t, txs)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, transactions, txs)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}
