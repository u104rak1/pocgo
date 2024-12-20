package authentication_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestSigninUsecase(t *testing.T) {
	type Mocks struct {
		authServ *domainMock.MockIAuthenticationService
		jwtServ  *appMock.MockIJWTService
	}

	var (
		userID      = idVO.NewUserIDForTest("user")
		email       = "sato@example.com"
		password    = "password"
		accessToken = "token"
		arg         = gomock.Any()
	)

	happyCmd := authApp.SigninCommand{
		Email:    email,
		Password: password,
	}

	tests := []struct {
		caseName string
		cmd      authApp.SigninCommand
		prepare  func(mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Positive: サインインが成功する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.authServ.EXPECT().Authenticate(arg, arg, arg).Return(&userID, nil)
				mocks.jwtServ.EXPECT().GenerateAccessToken(arg).Return(accessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: 認証に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.authServ.EXPECT().Authenticate(arg, arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: アクセストークンの生成に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.authServ.EXPECT().Authenticate(arg, arg, arg).Return(&userID, nil)
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
				authServ: domainMock.NewMockIAuthenticationService(ctrl),
				jwtServ:  appMock.NewMockIJWTService(ctrl),
			}
			uc := authApp.NewSigninUsecase(mocks.authServ, mocks.jwtServ)
			ctx := context.Background()
			tt.prepare(mocks)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, accessToken, dto.AccessToken)
			}
		})
	}
}
