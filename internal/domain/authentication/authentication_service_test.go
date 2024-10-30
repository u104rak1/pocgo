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
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestVerifyUniqueness(t *testing.T) {
	var (
		unknownErr = errors.New("unknown error")
	)

	tests := []struct {
		caseName string
		userID   string
		setup    func(ctx context.Context, mockAuthRepo *mock.MockIAuthenticationRepository, userID string)
		wantErr  error
	}{
		{
			caseName: "Successfully verifies that the user is unique.",
			userID:   ulid.GenerateStaticULID("Unique"),
			setup: func(ctx context.Context, mockAuthRepo *mock.MockIAuthenticationRepository, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs when the user already exists.",
			userID:   ulid.GenerateStaticULID("duplicate"),
			setup: func(ctx context.Context, mockAuthRepo *mock.MockIAuthenticationRepository, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(true, nil)
			},
			wantErr: authentication.ErrAlreadyExists,
		},
		{
			caseName: "Unknown Error occurs in ExistsByUserID.",
			userID:   ulid.GenerateStaticULID("unknown"),
			setup: func(ctx context.Context, mockAuthRepo *mock.MockIAuthenticationRepository, userID string) {
				mockAuthRepo.EXPECT().ExistsByUserID(ctx, userID).Return(false, unknownErr)
			},
			wantErr: unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
			mockUserRepo := mock.NewMockIUserRepository(ctrl)
			service := authentication.NewService(mockAuthRepo, mockUserRepo)
			ctx := context.Background()
			tt.setup(ctx, mockAuthRepo, tt.userID)

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
	accessToken, err := service.GenerateAccessToken(ctx, userID, jwtSecretKey)
	assert.NoError(t, err)

	tests := []struct {
		caseName     string
		accessToken  string
		jwtSecretKey []byte
		wantUserID   string
		wantErr      error
	}{
		{
			caseName:     "Successfully returns userID from valid access token.",
			accessToken:  accessToken,
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

func TestAuthenticate(t *testing.T) {
	type Mocks struct {
		authRepo *mock.MockIAuthenticationRepository
		userRepo *mock.MockIUserRepository
	}

	var (
		userID     = ulid.GenerateStaticULID("user")
		email      = "sato@examle.com"
		password   = "password"
		unknownErr = errors.New("unknown error")
	)
	user, err := userDomain.New(userID, "sato taro", email)
	assert.NoError(t, err)
	auth, err := authentication.New(userID, password)
	assert.NoError(t, err)

	tests := []struct {
		caseName   string
		email      string
		password   string
		setup      func(ctx context.Context, mocks Mocks)
		wantUserID string
		wantErr    error
	}{
		{
			caseName: "Successfully authenticates with valid email and password.",
			email:    email,
			password: password,
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(ctx, userID).Return(auth, nil)
			},
			wantUserID: userID,
			wantErr:    nil,
		},
		{
			caseName: "Error occurs when user is not found by email.",
			email:    "not-found@example.com",
			password: password,
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, "not-found@example.com").Return(nil, userDomain.ErrNotFound)
			},
			wantUserID: "",
			wantErr:    authentication.ErrAuthenticationFailed,
		},
		{
			caseName: "Unknown error occurs in FindByEmail.",
			email:    email,
			password: password,
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, email).Return(nil, unknownErr)
			},
			wantUserID: "",
			wantErr:    unknownErr,
		},
		{
			caseName: "Error occurs when authentication is not found by userID.",
			email:    email,
			password: password,
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(ctx, userID).Return(nil, authentication.ErrNotFound)
			},
			wantUserID: "",
			wantErr:    authentication.ErrAuthenticationFailed,
		},
		{
			caseName: "Unknown error occurs in FindByUserID.",
			email:    email,
			password: password,
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(ctx, userID).Return(nil, unknownErr)
			},
			wantUserID: "",
			wantErr:    unknownErr,
		},
		{
			caseName: "Error occurs when password is incorrect.",
			email:    email,
			password: "wrongPassword",
			setup: func(ctx context.Context, mocks Mocks) {
				mocks.userRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(ctx, userID).Return(auth, nil)
			},
			wantUserID: "",
			wantErr:    authentication.ErrAuthenticationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				authRepo: mock.NewMockIAuthenticationRepository(ctrl),
				userRepo: mock.NewMockIUserRepository(ctrl),
			}
			service := authentication.NewService(mocks.authRepo, mocks.userRepo)
			ctx := context.Background()
			tt.setup(ctx, mocks)

			userID, err := service.Authenticate(ctx, tt.email, tt.password)

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
