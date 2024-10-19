package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/authentication"
	"github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestVerifyUniqueness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)

	tests := []struct {
		caseName string
		userID   string
		setup    func(ctx context.Context, userID string)
		wantErr  error
	}{
		{
			caseName: "Happy path: return nil, if the user is unique.",
			userID:   ulid.GenerateStaticULID("Unique"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Edge case: return error, if the user is already exists.",
			userID:   ulid.GenerateStaticULID("duplicate"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(true, nil)
			},
			wantErr: authentication.ErrAuthenticationAlreadyExists,
		},
		{
			caseName: "Edge case: return error, if the repository returns an error.",
			userID:   ulid.GenerateStaticULID("error"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, errors.New("repository error"))
			},
			wantErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			service := authentication.NewService(mockAuthRepo)
			ctx := context.Background()
			tt.setup(ctx, tt.userID)
			err := service.VerifyUniqueness(ctx, tt.userID)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
