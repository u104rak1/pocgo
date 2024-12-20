package authentication_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
)

func TestSignupUsecase(t *testing.T) {
	type Mocks struct {
		userRepo *domainMock.MockIUserRepository
		userServ *domainMock.MockIUserService
		authRepo *domainMock.MockIAuthenticationRepository
		authServ *domainMock.MockIAuthenticationService
		jwtServ  *appMock.MockIJWTService
	}

	var (
		userName     = "sato taro"
		userEmail    = "sato@example.com"
		userPassword = "password"
		accessToken  = "token"
		arg          = gomock.Any()
	)

	happyCmd := authApp.SignupCommand{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}

	tests := []struct {
		caseName string
		cmd      authApp.SignupCommand
		prepare  func(mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Positive: サインアップが成功する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
				mocks.userRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(arg, arg).Return(nil)
				mocks.authRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.jwtServ.EXPECT().GenerateAccessToken(arg).Return(accessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: メールアドレスの一意性検証に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: ユーザー作成に失敗する",
			cmd:      authApp.SignupCommand{},
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: ユーザー保存に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
				mocks.userRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 認証の一意性検証に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
				mocks.userRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 認証保存に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
				mocks.userRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(arg, arg).Return(nil)
				mocks.authRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: アクセストークン生成に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().VerifyEmailUniqueness(arg, arg).Return(nil)
				mocks.userRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.authServ.EXPECT().VerifyUniqueness(arg, arg).Return(nil)
				mocks.authRepo.EXPECT().Save(arg, arg).Return(nil)
				mocks.jwtServ.EXPECT().GenerateAccessToken(arg).Return("", assert.AnError)
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
				jwtServ:  appMock.NewMockIJWTService(ctrl),
			}

			uc := authApp.NewSignupUsecase(mocks.userRepo, mocks.authRepo, mocks.userServ, mocks.authServ, mocks.jwtServ)
			ctx := context.Background()
			tt.prepare(mocks)

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
