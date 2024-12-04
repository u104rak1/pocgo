package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/u104raki/pocgo/internal/domain/mock"
	"github.com/u104raki/pocgo/internal/domain/user"
	"github.com/u104raki/pocgo/pkg/ulid"
)

func TestVerifyEmailUniqueness(t *testing.T) {
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
			wantErr: user.ErrEmailAlreadyExists,
		},
		{
			caseName: "Un known error occurs in ExsitsByEmail.",
			email:    "error@example.com",
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, email string) {
				mockUserRepo.EXPECT().ExistsByEmail(ctx, email).Return(false, assert.AnError)
			},
			wantErr: assert.AnError,
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

func TestEnsureUserExists(t *testing.T) {
	tests := []struct {
		caseName string
		id       string
		setup    func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, id string)
		wantErr  error
	}{
		{
			caseName: "Successfully verifies that the user exists.",
			id:       ulid.GenerateStaticULID("existing-user-id"),
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, id string) {
				mockUserRepo.EXPECT().ExistsByID(ctx, id).Return(true, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs when the user does not exist.",
			id:       ulid.GenerateStaticULID("non-existing-user-id"),
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, id string) {
				mockUserRepo.EXPECT().ExistsByID(ctx, id).Return(false, nil)
			},
			wantErr: user.ErrNotFound,
		},
		{
			caseName: "Unknown error occurs in ExistsByID.",
			id:       ulid.GenerateStaticULID("unknown-error-user-id"),
			setup: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository, id string) {
				mockUserRepo.EXPECT().ExistsByID(ctx, id).Return(false, assert.AnError)
			},
			wantErr: assert.AnError,
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
			tt.setup(ctx, mockUserRepo, tt.id)

			err := service.EnsureUserExists(ctx, tt.id)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
