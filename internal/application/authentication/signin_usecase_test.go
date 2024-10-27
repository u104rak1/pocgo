package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	authApp "github.com/ucho456job/pocgo/internal/application/authentication"
	authDomain "github.com/ucho456job/pocgo/internal/domain/authentication"
	domainMock "github.com/ucho456job/pocgo/internal/domain/mock"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSigninUsecase(t *testing.T) {
	type Mocks struct {
		mockUserRepo *domainMock.MockIUserRepository
		mockAuthRepo *domainMock.MockIAuthenticationRepository
		mockAuthServ *domainMock.MockIAuthenticationService
	}

	var (
		validEmail       = "sato@example.com"
		validPassword    = "password"
		validAccessToken = "token"
	)

	userID := ulid.GenerateStaticULID("user")
	user, err := userDomain.New(userID, "sato", validEmail)
	assert.NoError(t, err)
	auth, err := authDomain.New(userID, validPassword)
	assert.NoError(t, err)

	validCmd := authApp.SigninCommand{
		Email:    validEmail,
		Password: validPassword,
	}

	tests := []struct {
		caseName string
		cmd      authApp.SigninCommand
		prepare  func(ctx context.Context, mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Signin is successfully done.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserRepo.EXPECT().FindByEmail(ctx, validEmail).Return(user, nil)
				mocks.mockAuthRepo.EXPECT().FindByUserID(ctx, userID).Return(auth, nil)
				mocks.mockAuthServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return(validAccessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs because the user does not exist.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserRepo.EXPECT().FindByEmail(ctx, validEmail).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs because the authentication does not exist.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserRepo.EXPECT().FindByEmail(ctx, validEmail).Return(user, nil)
				mocks.mockAuthRepo.EXPECT().FindByUserID(ctx, userID).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs because the password is incorrect.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				auth, err := authDomain.New(userID, "invalidPassword")
				assert.NoError(t, err)
				mocks.mockUserRepo.EXPECT().FindByEmail(ctx, validEmail).Return(user, nil)
				mocks.mockAuthRepo.EXPECT().FindByUserID(ctx, userID).Return(auth, nil)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during access token generation.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserRepo.EXPECT().FindByEmail(ctx, validEmail).Return(user, nil)
				mocks.mockAuthRepo.EXPECT().FindByUserID(ctx, userID).Return(auth, nil)
				mocks.mockAuthServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return("", errors.New("error"))
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
				mockUserRepo: domainMock.NewMockIUserRepository(ctrl),
				mockAuthRepo: domainMock.NewMockIAuthenticationRepository(ctrl),
				mockAuthServ: domainMock.NewMockIAuthenticationService(ctrl),
			}

			uc := authApp.NewSigninUsecase(mocks.mockUserRepo, mocks.mockAuthRepo, mocks.mockAuthServ)
			ctx := context.Background()
			tt.prepare(ctx, mocks)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validAccessToken, dto.AccessToken)
			}
		})
	}
}
