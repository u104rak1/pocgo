package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountApp "github.com/ucho456job/pocgo/internal/application/account"
	"github.com/ucho456job/pocgo/internal/application/authentication"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	domainMock "github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestSignupUsecase(t *testing.T) {
	type Mocks struct {
		userUC    *appMock.MockICreateUserUsecase
		accountUC *appMock.MockICreateAccountUsecase
		authServ  *domainMock.MockIAuthenticationService
	}

	var (
		userID           = ulid.GenerateStaticULID("user")
		userName         = "sato taro"
		userEmail        = "sato@example.com"
		userPassword     = "password"
		accountID        = ulid.GenerateStaticULID("account")
		accountName      = "For work"
		accountBalance   = 0.0
		accountPassword  = "1234"
		accountCurrency  = money.JPY
		accountUpdatedAt = timer.Now().String()
		accessToken      = "token"
	)

	cmd := authentication.SignupCommand{
		User: userApp.CreateUserCommand{
			Name:     userName,
			Email:    userEmail,
			Password: userPassword,
		},
		Account: accountApp.CreateAccountCommand{
			Name:     accountName,
			Password: accountPassword,
			Currency: accountCurrency,
		},
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
				mocks.userUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				}, nil)
				mocks.accountUC.EXPECT().Run(ctx, gomock.Any()).Return(&accountApp.CreateAccountDTO{
					ID:        accountID,
					UserID:    userID,
					Name:      accountName,
					Balance:   accountBalance,
					Currency:  accountCurrency,
					UpdatedAt: accountUpdatedAt,
				}, nil)
				mocks.authServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return(accessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during userCreateDto creation.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during accountCreateDto creation.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				}, nil)
				mocks.accountUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during access token generation.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				}, nil)
				mocks.accountUC.EXPECT().Run(ctx, gomock.Any()).Return(&accountApp.CreateAccountDTO{
					ID:        accountID,
					UserID:    userID,
					Name:      accountName,
					Balance:   accountBalance,
					Currency:  accountCurrency,
					UpdatedAt: accountUpdatedAt,
				}, nil)
				mocks.authServ.EXPECT().GenerateAccessToken(ctx, userID, gomock.Any()).Return("", errors.New("error"))
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
				userUC:    appMock.NewMockICreateUserUsecase(ctrl),
				accountUC: appMock.NewMockICreateAccountUsecase(ctrl),
				authServ:  domainMock.NewMockIAuthenticationService(ctrl),
			}
			mockUnitOfWork := &appMock.MockIUnitOfWorkWithResult[authentication.SignupDTO]{}

			uc := authentication.NewSignupUsecase(mocks.userUC, mocks.accountUC, mocks.authServ, mockUnitOfWork)
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
				assert.Equal(t, tt.cmd.User.Name, dto.User.Name)
				assert.Equal(t, tt.cmd.User.Email, dto.User.Email)
				assert.NotEmpty(t, dto.Account.ID)
				assert.Equal(t, dto.User.ID, dto.Account.UserID)
				assert.Equal(t, tt.cmd.Account.Name, dto.Account.Name)
				assert.Equal(t, 0.0, dto.Account.Balance)
				assert.Equal(t, tt.cmd.Account.Currency, dto.Account.Currency)
				assert.NotEmpty(t, dto.Account.UpdatedAt)
				assert.Equal(t, accessToken, dto.AccessToken)
			}
		})
	}
}
