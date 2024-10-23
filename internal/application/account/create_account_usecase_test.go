package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	accountUC "github.com/ucho456job/pocgo/internal/application/account"
	"github.com/ucho456job/pocgo/internal/domain/mock"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestCreateAccountUsecase(t *testing.T) {
	userID := ulid.GenerateStaticULID("user")
	validCmd := accountUC.CreateAccountCommand{
		UserID:   userID,
		Name:     "For work",
		Password: "1234",
		Currency: money.JPY,
	}

	tests := []struct {
		caseName string
		cmd      accountUC.CreateAccountCommand
		prepare  func(ctx context.Context, mockRepo *mock.MockIAccountRepository)
		wantErr  bool
	}{
		{
			caseName: "Account is successfully created.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mockRepo *mock.MockIAccountRepository) {
				mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during accountDomain creation.",
			cmd:      accountUC.CreateAccountCommand{},
			prepare:  func(ctx context.Context, mockRepo *mock.MockIAccountRepository) {},
			wantErr:  true,
		},
		{
			caseName: "Error occurs during Save in accountRepository.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mockRepo *mock.MockIAccountRepository) {
				mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockIAccountRepository(ctrl)
			uc := accountUC.NewCreateAccountUsecase(mockRepo)
			ctx := context.Background()
			tt.prepare(ctx, mockRepo)

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
