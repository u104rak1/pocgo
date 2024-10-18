package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/domain/user/mock"
)

func TestVerifyEmailUniquenessService_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockIUserRepository(ctrl)

	tests := []struct {
		name    string
		email   string
		setup   func(ctx context.Context, email string)
		wantErr error
	}{
		{
			name:  "Happy path: メールがUniqueの時、エラーが発生しない",
			email: "new@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			name:  "Edge case: メールが重複している時、エラーが発生する",
			email: "existing@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(true, nil)
			},
			wantErr: user.ErrUserEmailAlreadyExists,
		},
		{
			name:  "Edge case: ExistsByEmailでエラーが発生する",
			email: "error@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, errors.New("repository error"))
			},
			wantErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := user.NewVerifyEmailUniquenessService(mockUserRepo)
			ctx := context.Background()
			tt.setup(ctx, tt.email)
			err := service.Run(ctx, tt.email)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
