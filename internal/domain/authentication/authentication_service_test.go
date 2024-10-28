package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
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
	mockUserRepo := mock.NewMockIUserRepository(ctrl)
	errRepo := errors.New("repository error")

	tests := []struct {
		caseName string
		userID   string
		setup    func(ctx context.Context, userID string)
		wantErr  error
	}{
		{
			caseName: "Successfully verifies that the user is unique.",
			userID:   ulid.GenerateStaticULID("Unique"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs when the user already exists.",
			userID:   ulid.GenerateStaticULID("duplicate"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(true, nil)
			},
			wantErr: authentication.ErrAuthenticationAlreadyExists,
		},
		{
			caseName: "Error occurs when the repository returns an error.",
			userID:   ulid.GenerateStaticULID("error"),
			setup: func(ctx context.Context, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			service := authentication.NewService(mockAuthRepo, mockUserRepo)
			ctx := context.Background()
			tt.setup(ctx, tt.userID)
			err := service.VerifyUniqueness(ctx, tt.userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	t.Run("Successfully returns access token.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
		mockUserRepo := mock.NewMockIUserRepository(ctrl)
		service := authentication.NewService(mockAuthRepo, mockUserRepo)

		ctx := context.Background()
		userID := ulid.GenerateStaticULID("user")
		jwtSecretKey := []byte("validSecretKey")
		token, err := service.GenerateAccessToken(ctx, userID, jwtSecretKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestGetUserIDFromAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
	mockUserRepo := mock.NewMockIUserRepository(ctrl)
	service := authentication.NewService(mockAuthRepo, mockUserRepo)

	ctx := context.Background()
	userID := ulid.GenerateStaticULID("user")
	jwtSecretKey := []byte("validSecretKey")
	validAccessToken, _ := service.GenerateAccessToken(ctx, userID, jwtSecretKey)

	tests := []struct {
		caseName     string
		accessToken  string
		jwtSecretKey []byte
		wantUserID   string
		wantErr      error
	}{
		{
			caseName:     "Successfully returns userID from valid access token.",
			accessToken:  validAccessToken,
			jwtSecretKey: jwtSecretKey,
			wantUserID:   userID,
			wantErr:      nil,
		},
		{
			caseName:     "Error occurs when the token is invalid.",
			accessToken:  "invalidToken",
			jwtSecretKey: jwtSecretKey,
			wantUserID:   "",
			wantErr:      jwt.ErrTokenMalformed,
		},
		// All branch testing is difficult, so skip it.
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			userID, err := service.GetUserIDFromAccessToken(ctx, tt.accessToken, tt.jwtSecretKey)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUserID, userID)
			}
		})
	}
}
