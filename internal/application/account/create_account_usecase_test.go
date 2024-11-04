package account_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountUC "github.com/ucho456job/pocgo/internal/application/account"
	appMock "github.com/ucho456job/pocgo/internal/application/mock"
	"github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestCreateAccountUsecase(t *testing.T) {
	type Mocks struct {
		accountRepo *mock.MockIAccountRepository
		accountServ *mock.MockIAccountService
		userServ    *mock.MockIUserService
	}

	userID := ulid.GenerateStaticULID("user")
	cmd := accountUC.CreateAccountCommand{
		UserID:   userID,
		Name:     "For work",
		Password: "1234",
		Currency: money.JPY,
	}

	tests := []struct {
		caseName string
		cmd      accountUC.CreateAccountCommand
		prepare  func(ctx context.Context, mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Account is successfully created.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(ctx, userID).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(ctx, userID).Return(nil)
				mocks.accountRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			caseName: "An error occurs during account creation.",
			cmd:      accountUC.CreateAccountCommand{},
			prepare:  func(ctx context.Context, mocks Mocks) {},
			wantErr:  true,
		},
		{
			caseName: "An error occurs when the user not found.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(ctx, userID).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "An error occurs if the user has reached the limit of accounts.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(ctx, userID).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(ctx, userID).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "An error occurs during Save in accountRepository.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(ctx, userID).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(ctx, userID).Return(nil)
				mocks.accountRepo.EXPECT().Save(ctx, gomock.Any()).Return(assert.AnError)
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
				accountRepo: mock.NewMockIAccountRepository(ctrl),
				accountServ: mock.NewMockIAccountService(ctrl),
				userServ:    mock.NewMockIUserService(ctrl),
			}
			mockUnitOfWork := &appMock.MockIUnitOfWork{}

			uc := accountUC.NewCreateAccountUsecase(
				mocks.accountRepo, mocks.accountServ, mocks.userServ, mockUnitOfWork,
			)
			ctx := context.Background()
			tt.prepare(ctx, mocks)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.NotEmpty(t, dto.ID)
				assert.Equal(t, tt.cmd.UserID, dto.UserID)
				assert.Equal(t, tt.cmd.Name, dto.Name)
				assert.Equal(t, 0.0, dto.Balance)
				assert.Equal(t, tt.cmd.Currency, dto.Currency)
				assert.NotEmpty(t, dto.UpdatedAt)
			}
		})
	}
}
