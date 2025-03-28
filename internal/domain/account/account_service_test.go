package account_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/pkg/strutil"
)

func TestCheckLimit(t *testing.T) {
	var (
		userID = idVO.NewUserIDForTest("user")
		arg    = gomock.Any()
	)

	tests := []struct {
		caseName string
		userID   idVO.UserID
		setup    func(mockAccountRepo *mock.MockIAccountRepository)
		errMsg   string
	}{
		{
			caseName: "Positive: 口座数が上限に達していない場合はエラーが返らない",
			userID:   userID,
			setup: func(mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(arg, arg).Return(2, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: 口座数が上限に達している場合はエラーが返る",
			userID:   userID,
			setup: func(mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(arg, arg).Return(3, nil)
			},
			errMsg: "account limit reached, maximum 3 accounts",
		},
		{
			caseName: "Negative: CountByUserIDでエラーが返る場合はエラーが返る",
			userID:   userID,
			setup: func(mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(arg, arg).Return(0, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockIAccountRepository(ctrl)
			service := accountDomain.NewService(mockAccountRepo)
			ctx := context.Background()
			tt.setup(mockAccountRepo)

			err := service.CheckLimit(ctx, tt.userID)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAndAuthorize(t *testing.T) {
	var (
		accountID = idVO.NewAccountIDForTest("account")
		userID    = idVO.NewUserIDForTest("user")
		name      = "account-name"
		password  = "1234"
		amount    = 100.0
		currency  = moneyVO.JPY
		arg       = gomock.Any()
	)

	tests := []struct {
		caseName  string
		accountID idVO.AccountID
		userID    *idVO.UserID
		password  *string
		setup     func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account)
		errMsg    string
	}{
		{
			caseName:  "Positive: ユーザーIDとパスワードが一致する場合は口座が取得できる",
			accountID: accountID,
			userID:    &userID,
			password:  &password,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "Positive: ユーザーIDがnilの場合は、ユーザーIDの検証を行わないで口座が取得できる",
			accountID: accountID,
			userID:    nil,
			password:  &password,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "Positive: パスワードがnilの場合は、パスワードの検証を行わないで口座が取得できる",
			accountID: accountID,
			userID:    &userID,
			password:  nil,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "Negative: FindByIDでエラーが返る場合はエラーが返る",
			accountID: accountID,
			userID:    &userID,
			password:  nil,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(nil, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName:  "Negative: 口座が存在しない場合はエラーが返る",
			accountID: accountID,
			userID:    nil,
			password:  nil,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(nil, nil)
			},
			errMsg: "account not found",
		},
		{
			caseName:  "Negative: ユーザーIDが一致しない場合はエラーが返る",
			accountID: accountID,
			userID: func() *idVO.UserID {
				id := idVO.NewUserIDForTest("unauthorized-user")
				return &id
			}(),
			password: nil,
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(account, nil)
			},
			errMsg: "unauthorized access to account",
		},
		{
			caseName:  "Negative: パスワードが一致しない場合はエラーが返る",
			accountID: accountID,
			userID:    nil,
			password:  strutil.StrPointer("5678"),
			setup: func(mockAccountRepo *mock.MockIAccountRepository, account *accountDomain.Account) {
				mockAccountRepo.EXPECT().FindByID(arg, arg).Return(account, nil)
			},
			errMsg: "passwords do not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockIAccountRepository(ctrl)
			service := accountDomain.NewService(mockAccountRepo)
			ctx := context.Background()
			account, err := accountDomain.New(userID, amount, name, password, currency)
			assert.NoError(t, err)
			tt.setup(mockAccountRepo, account)

			a, err := service.GetAndAuthorize(ctx, tt.accountID, tt.userID, tt.password)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, account.ID(), a.ID())
			}
		})
	}
}
