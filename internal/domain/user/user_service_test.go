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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockIUserRepository(ctrl)

	tests := []struct {
		caseName string
		email    string
		setup    func(ctx context.Context, email string)
		wantErr  error
	}{
		{
			caseName: "Happy path: return nil, if the email is unique.",
			email:    "new@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Edge case: return error, if the email is already exists.",
			email:    "existing@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(true, nil)
			},
			wantErr: user.ErrUserEmailAlreadyExists,
		},
		{
			caseName: "Edge case: return error, if the repository returns an error.",
			email:    "error@example.com",
			setup: func(ctx context.Context, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, errors.New("repository error"))
			},
			wantErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			service := user.NewService(mockUserRepo)
			ctx := context.Background()
			tt.setup(ctx, tt.email)
			err := service.VerifyEmailUniqueness(ctx, tt.email)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}