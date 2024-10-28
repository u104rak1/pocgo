package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	domainMock "github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSigninUsecase(t *testing.T) {
	var (
		userID      = ulid.GenerateStaticULID("user")
		email       = "sato@example.com"
		password    = "password"
		accessToken = "token"
	)

	cmd := authApp.SigninCommand{
		Email:    email,
		Password: password,
	}

	tests := []struct {
		caseName string
		cmd      authApp.SigninCommand
		prepare  func(ctx context.Context, authServ *domainMock.MockIAuthenticationService)
		wantErr  bool
	}{
		{
			caseName: "Signin is successfully done.",
			cmd:      cmd,
			prepare: func(ctx context.Context, authServ *domainMock.MockIAuthenticationService) {
				authServ.EXPECT().Authenticate(ctx, email, password).Return(userID, nil)
				authServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return(accessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs when authentication fails.",
			cmd:      cmd,
			prepare: func(ctx context.Context, authServ *domainMock.MockIAuthenticationService) {
				authServ.EXPECT().Authenticate(ctx, email, password).Return("", errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs when generating an access token fails.",
			cmd:      cmd,
			prepare: func(ctx context.Context, authServ *domainMock.MockIAuthenticationService) {
				authServ.EXPECT().Authenticate(ctx, email, password).Return(userID, nil)
				authServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return("", errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authServ := domainMock.NewMockIAuthenticationService(ctrl)
			uc := authApp.NewSigninUsecase(authServ)
			ctx := context.Background()
			tt.prepare(ctx, authServ)

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
