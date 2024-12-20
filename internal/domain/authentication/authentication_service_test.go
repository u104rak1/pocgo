package authentication_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestVerifyUniqueness(t *testing.T) {
	arg := gomock.Any()

	tests := []struct {
		caseName string
		userID   idVO.UserID
		setup    func(mockAuthRepo *mock.MockIAuthenticationRepository)
		errMsg   string
	}{
		{
			caseName: "Positive: ユーザーがユニークな場合はエラーが返らない",
			userID:   idVO.NewUserIDForTest("Unique"),
			setup: func(mockAuthRepo *mock.MockIAuthenticationRepository) {
				mockAuthRepo.EXPECT().ExistsByUserID(arg, arg).Return(false, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: ユーザーが既に存在する場合はエラーが返る",
			userID:   idVO.NewUserIDForTest("duplicate"),
			setup: func(mockAuthRepo *mock.MockIAuthenticationRepository) {
				mockAuthRepo.EXPECT().ExistsByUserID(arg, arg).Return(true, nil)
			},
			errMsg: "authentication already exists",
		},
		{
			caseName: "Negative: ExistsByUserIDでエラーが返る場合はエラーが返る",
			userID:   idVO.NewUserIDForTest("unknown"),
			setup: func(mockAuthRepo *mock.MockIAuthenticationRepository) {
				mockAuthRepo.EXPECT().ExistsByUserID(arg, arg).Return(false, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
			mockUserRepo := mock.NewMockIUserRepository(ctrl)
			service := authDomain.NewService(mockAuthRepo, mockUserRepo)
			ctx := context.Background()
			tt.setup(mockAuthRepo)

			err := service.VerifyUniqueness(ctx, tt.userID)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
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
		userID   = idVO.NewUserIDForTest("user")
		name     = "sato taro"
		email    = "sato@examle.com"
		password = "password"
		arg      = gomock.Any()
	)

	tests := []struct {
		caseName   string
		email      string
		password   string
		setup      func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication)
		wantUserID idVO.UserID
		errMsg     string
	}{
		{
			caseName: "Positive: 有効なメールアドレスとパスワードの場合は認証できる",
			email:    email,
			password: password,
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(arg, arg).Return(auth, nil)
			},
			wantUserID: userID,
			errMsg:     "",
		},
		{
			caseName: "Negative: ユーザーが見つからない場合はエラーが返る",
			email:    "not-found@example.com",
			password: password,
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(nil, nil)
			},
			errMsg: "email or password is incorrect",
		},
		{
			caseName: "Negative: FindByEmailが失敗する場合はエラーが返る",
			email:    email,
			password: password,
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(nil, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: 認証情報が見つからない場合はエラーが返る",
			email:    email,
			password: password,
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(arg, arg).Return(nil, nil)
			},
			errMsg: "email or password is incorrect",
		},
		{
			caseName: "Negative: FindByUserIDが失敗する場合はエラーが返る",
			email:    email,
			password: password,
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(arg, arg).Return(nil, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
		{
			caseName: "Negative: パスワードが一致しない場合はエラーが返る",
			email:    email,
			password: "wrongPassword",
			setup: func(mocks Mocks, user *userDomain.User, auth *authDomain.Authentication) {
				mocks.userRepo.EXPECT().FindByEmail(arg, arg).Return(user, nil)
				mocks.authRepo.EXPECT().FindByUserID(arg, arg).Return(auth, nil)
			},
			errMsg: "email or password is incorrect",
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
			service := authDomain.NewService(mocks.authRepo, mocks.userRepo)
			ctx := context.Background()
			user, err := userDomain.Reconstruct(userID.String(), name, email)
			assert.NoError(t, err)
			auth, err := authDomain.New(user.ID(), password)
			assert.NoError(t, err)
			tt.setup(mocks, user, auth)

			gotUserID, err := service.Authenticate(ctx, tt.email, tt.password)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Empty(t, gotUserID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUserID, *gotUserID)
			}
		})
	}
}
