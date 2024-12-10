package authentication_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/application/authentication"
	"github.com/u104rak1/pocgo/internal/config"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
)

func TestSignupUsecase(t *testing.T) {
	type Mocks struct {
		userRepo *domainMock.MockIUserRepository
		userServ *domainMock.MockIUserService
		authRepo *domainMock.MockIAuthenticationRepository
		authServ *domainMock.MockIAuthenticationService
	}

	var (
		userName     = "sato taro"
		userEmail    = "sato@example.com"
		userPassword = "password"
		accessToken  = "token"
		jwtSecretKey = []byte(config.NewEnv().JWT_SECRET_KEY)
	)

	cmd := authentication.SignupCommand{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}

	tests := []struct {
		caseName string
		cmd      authentication.SignupCommand
		prepare  func(ctx context.Context, mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Signup is successfully done.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(nil)
				mocks.userRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.authRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().GenerateAccessToken(ctx, gomock.Any(), jwtSecretKey).Return(accessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during email uniqueness verification.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during user save.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(nil)
				mocks.userRepo.EXPECT().Save(ctx, gomock.Any()).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during authentication uniqueness verification.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(nil)
				mocks.userRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during authentication save.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(nil)
				mocks.userRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.authRepo.EXPECT().Save(ctx, gomock.Any()).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during access token generation.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(ctx, userEmail).Return(nil)
				mocks.userRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.authRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.authServ.EXPECT().GenerateAccessToken(ctx, gomock.Any(), jwtSecretKey).Return("", assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := Mocks{
				userRepo: domainMock.NewMockIUserRepository(ctrl),
				userServ: domainMock.NewMockIUserService(ctrl),
				authRepo: domainMock.NewMockIAuthenticationRepository(ctrl),
				authServ: domainMock.NewMockIAuthenticationService(ctrl),
			}

			uc := authentication.NewSignupUsecase(mocks.userRepo, mocks.authRepo, mocks.userServ, mocks.authServ)
			ctx := context.Background()
			tt.prepare(ctx, mocks)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.NotEmpty(t, dto.User.ID)
				assert.Equal(t, tt.cmd.Name, dto.User.Name)
				assert.Equal(t, tt.cmd.Email, dto.User.Email)
				assert.Equal(t, accessToken, dto.AccessToken)
			}
		})
	}
}
