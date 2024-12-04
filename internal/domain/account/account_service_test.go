package account_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104raki/pocgo/internal/domain/account"
	"github.com/u104raki/pocgo/internal/domain/mock"
	"github.com/u104raki/pocgo/internal/domain/value_object/money"
	"github.com/u104raki/pocgo/pkg/timer"
	"github.com/u104raki/pocgo/pkg/ulid"
)

func TestCheckLimit(t *testing.T) {
	var userID = ulid.GenerateStaticULID("user")

	tests := []struct {
		caseName string
		userID   string
		setup    func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository)
		errMsg   string
	}{
		{
			caseName: "Successfully returns no error when account count is below the limit (count = 2).",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(2, nil)
			},
			errMsg: "",
		},
		{
			caseName: "The Limit Reached Error occurs when account count reaches the limit (count = 3).",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(3, nil)
			},
			errMsg: "account limit reached, maximum 3 accounts",
		},
		{
			caseName: "An unknown error occurs in CountByUserID.",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(0, assert.AnError)
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
			tt.setup(ctx, mockAccountRepo)

			err := service.CheckLimit(ctx, tt.userID)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGetAndAuthorize(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "account-name"
		password  = "1234"
		amount    = 100.0
		currency  = money.JPY
		updatedAt = timer.GetFixedDate()
	)

	account, err := accountDomain.New(accountID, userID, name, password, amount, currency, updatedAt)
	assert.NoError(t, err)

	tests := []struct {
		caseName  string
		accountID string
		userID    *string
		password  *string
		setup     func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository)
		errMsg    string
	}{
		{
			caseName:  "Successfully retrieves account when verify user id and password.",
			accountID: accountID,
			userID:    &userID,
			password:  &password,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "Successfully retrieves account without userID verification.",
			accountID: accountID,
			userID:    nil,
			password:  &password,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "Successfully retrieves account without password verification.",
			accountID: accountID,
			userID:    &userID,
			password:  nil,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(account, nil)
			},
			errMsg: "",
		},
		{
			caseName:  "An unknown error occurs in FindByID.",
			accountID: accountID,
			userID:    &userID,
			password:  nil,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(nil, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName:  "The account not found error occurs when the account does not exist.",
			accountID: accountID,
			userID:    nil,
			password:  nil,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(nil, nil)
			},
			errMsg: "account not found",
		},
		{
			caseName:  "The unauthorized access error occurs when the user ID does not match.",
			accountID: accountID,
			userID: func() *string {
				unauthorizedUserID := ulid.GenerateStaticULID("unauthorized-user")
				return &unauthorizedUserID
			}(),
			password: nil,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(account, nil)
			},
			errMsg: "unauthorized access to account",
		},
		{
			caseName:  "The unmatched password error occurs when the password does not match.",
			accountID: accountID,
			userID:    nil,
			password: func() *string {
				unmatchedPassword := "5678"
				return &unmatchedPassword
			}(),
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().FindByID(ctx, accountID).Return(account, nil)
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
			tt.setup(ctx, mockAccountRepo)

			a, err := service.GetAndAuthorize(ctx, tt.accountID, tt.userID, tt.password)
			if tt.errMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, account.ID(), a.ID())
			} else {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errMsg)
			}
		})
	}
}
