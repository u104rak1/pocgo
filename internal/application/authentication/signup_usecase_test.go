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

func TestSignupUsecase_Run(t *testing.T) {
	type Mocks struct {
		mockUserUC    *appMock.MockICreateUserUsecase
		mockAccountUC *appMock.MockICreateAccountUsecase
		mockAuthServ  *domainMock.MockIAuthenticationService
	}

	var (
		validUserID           = ulid.GenerateStaticULID("user")
		validUserName         = "sato taro"
		validUserEmail        = "sato@example.com"
		validUserPassword     = "password"
		validAccountID        = ulid.GenerateStaticULID("account")
		validAccountName      = "For work"
		validAccountBalance   = 0.0
		validAccountPassword  = "1234"
		validAccountCurrency  = money.JPY
		validAccountUpdatedAt = timer.Now().String()
		validAccessToken      = "token"
		err                   = errors.New("error")
	)

	validCmd := authentication.SignupCommand{
		User: userApp.CreateUserCommand{
			Name:     validUserName,
			Email:    validUserEmail,
			Password: validUserPassword,
		},
		Account: accountApp.CreateAccountCommand{
			Name:     validAccountName,
			Password: validAccountPassword,
			Currency: validAccountCurrency,
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
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    validUserID,
					Name:  validUserName,
					Email: validUserEmail,
				}, nil)
				mocks.mockAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(&accountApp.CreateAccountDTO{
					ID:        validAccountID,
					UserID:    validUserID,
					Name:      validAccountName,
					Balance:   validAccountBalance,
					Currency:  validAccountCurrency,
					UpdatedAt: validAccountUpdatedAt,
				}, nil)
				mocks.mockAuthServ.EXPECT().GenerateAccessToken(ctx, validUserID, gomock.Any()).Return(validAccessToken, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during userCreateDto creation.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, err)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during accountCreateDto creation.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    validUserID,
					Name:  validUserName,
					Email: validUserEmail,
				}, nil)
				mocks.mockAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(nil, err)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during access token generation.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserUC.EXPECT().Run(ctx, gomock.Any()).Return(&userApp.CreateUserDTO{
					ID:    validUserID,
					Name:  validUserName,
					Email: validUserEmail,
				}, nil)
				mocks.mockAccountUC.EXPECT().Run(ctx, gomock.Any()).Return(&accountApp.CreateAccountDTO{
					ID:        validAccountID,
					UserID:    validUserID,
					Name:      validAccountName,
					Balance:   validAccountBalance,
					Currency:  validAccountCurrency,
					UpdatedAt: validAccountUpdatedAt,
				}, nil)
				mocks.mockAuthServ.EXPECT().GenerateAccessToken(ctx, validUserID, gomock.Any()).Return("", err)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUC := appMock.NewMockICreateUserUsecase(ctrl)
			mockAccountUC := appMock.NewMockICreateAccountUsecase(ctrl)
			mockAuthServ := domainMock.NewMockIAuthenticationService(ctrl)
			mockUnitOfWork := &appMock.MockIUnitOfWorkWithResult[authentication.SignupDTO]{}
			uc := authentication.NewSignupUsecase(mockUserUC, mockAccountUC, mockAuthServ, mockUnitOfWork)
			ctx := context.Background()
			mocks := Mocks{
				mockUserUC:    mockUserUC,
				mockAccountUC: mockAccountUC,
				mockAuthServ:  mockAuthServ,
			}
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
				assert.Equal(t, validAccessToken, dto.AccessToken)
			}
		})
	}
}
