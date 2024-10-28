package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/internal/domain/user"
)

func TestVerifyEmailUniqueness(t *testing.T) {
	var unknownErr = errors.New("unknown error")

	tests := []struct {
		caseName string
		email    string
		setup    func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, email string)
		wantErr  error
	}{
		{
			caseName: "Successfully verifies that the email is unique.",
			email:    "new@example.com",
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs when the email already exists.",
			email:    "existing@example.com",
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(true, nil)
			},
			wantErr: user.ErrUserEmailAlreadyExists,
		},
		{
			caseName: "Un known error occurs in ExsitsByEmail.",
			email:    "error@example.com",
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, unknownErr)
			},
			wantErr: unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock.NewMockIUserRepository(ctrl)

			service := user.NewService(mockUserRepo)
			ctx := context.Background()
			tt.setup(ctx, mockUserRepo, tt.email)

			err := service.VerifyEmailUniqueness(ctx, tt.email)

			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
