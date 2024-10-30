package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/account"
	"github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestCheckLimit(t *testing.T) {
	var (
		userID     = ulid.GenerateStaticULID("user")
		unknownErr = errors.New("unknown error")
	)

	tests := []struct {
		caseName string
		userID   string
		setup    func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository)
		wantErr  error
	}{
		{
			caseName: "Successfully returns no error when account count is below the limit (count = 2).",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(2, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs when account count reaches the limit (count = 3).",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(3, nil)
			},
			wantErr: account.ErrLimitReached,
		},
		{
			caseName: "Error occurs when account count exceeds the limit (count = 4).",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(4, nil)
			},
			wantErr: account.ErrLimitReached,
		},
		{
			caseName: "Unknown error occurs in CountByUserID.",
			userID:   userID,
			setup: func(ctx context.Context, mockAccountRepo *mock.MockIAccountRepository) {
				mockAccountRepo.EXPECT().CountByUserID(ctx, userID).Return(0, unknownErr)
			},
			wantErr: unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mock.NewMockIAccountRepository(ctrl)

			service := account.NewService(mockAccountRepo)
			ctx := context.Background()
			tt.setup(ctx, mockAccountRepo)

			err := service.CheckLimit(ctx, tt.userID)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
