package account_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountUC "github.com/u104rak1/pocgo/internal/application/account"
	appMock "github.com/u104rak1/pocgo/internal/application/mock"
	domainMock "github.com/u104rak1/pocgo/internal/domain/mock"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
)

func TestCreateAccountUsecase(t *testing.T) {
	type Mocks struct {
		accountRepo *domainMock.MockIAccountRepository
		accountServ *domainMock.MockIAccountService
		userServ    *domainMock.MockIUserService
	}

	var (
		userID   = idVO.NewUserIDForTest("user")
		name     = "For work"
		password = "1234"
		currency = moneyVO.JPY
		arg      = gomock.Any()
	)

	happyCmd := accountUC.CreateAccountCommand{
		UserID:   userID.String(),
		Name:     name,
		Password: password,
		Currency: currency,
	}

	tests := []struct {
		caseName string
		cmd      accountUC.CreateAccountCommand
		prepare  func(mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "Positive: 口座作成が成功する",
			cmd: accountUC.CreateAccountCommand{
				UserID: "invalid",
			},
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(arg, arg).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(arg, arg).Return(nil)
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: ユーザーIDが不正な形式である",
			cmd: accountUC.CreateAccountCommand{
				UserID: "invalid",
			},
			prepare: func(mocks Mocks) {},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座作成に失敗する",
			cmd:      accountUC.CreateAccountCommand{},
			prepare:  func(mocks Mocks) {},
			wantErr:  true,
		},
		{
			caseName: "Negative: ユーザーの存在確認に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: ユーザーの所持口座上限確認に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(arg, arg).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(arg, arg).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: 口座保存に失敗する",
			cmd:      happyCmd,
			prepare: func(mocks Mocks) {
				mocks.userServ.EXPECT().EnsureUserExists(arg, arg).Return(nil)
				mocks.accountServ.EXPECT().CheckLimit(arg, arg).Return(nil)
				mocks.accountRepo.EXPECT().Save(arg, arg).Return(assert.AnError)
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
				accountRepo: domainMock.NewMockIAccountRepository(ctrl),
				accountServ: domainMock.NewMockIAccountService(ctrl),
				userServ:    domainMock.NewMockIUserService(ctrl),
			}
			mockUnitOfWork := &appMock.MockIUnitOfWork{}

			uc := accountUC.NewCreateAccountUsecase(
				mocks.accountRepo, mocks.accountServ, mocks.userServ, mockUnitOfWork,
			)
			ctx := context.Background()
			tt.prepare(mocks)

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
