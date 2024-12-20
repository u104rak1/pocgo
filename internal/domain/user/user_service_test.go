package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestVerifyEmailUniqueness(t *testing.T) {
	var arg = gomock.Any()

	tests := []struct {
		caseName string
		email    string
		setup    func(mockUserRepo *mock.MockIUserRepository)
		errMsg   string
	}{
		{
			caseName: "Positive: ユーザーのメールアドレスがユニークな場合はエラーが返らない",
			email:    "new@example.com",
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByEmail(arg, arg).Return(false, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: ユーザーのメールアドレスが既に存在する場合はエラーが返る",
			email:    "existing@example.com",
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByEmail(arg, arg).Return(true, nil)
			},
			errMsg: "user email already exists",
		},
		{
			caseName: "Negative: ExistsByEmailでエラーが返る場合はエラーが返る",
			email:    "error@example.com",
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByEmail(arg, arg).Return(false, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock.NewMockIUserRepository(ctrl)

			service := userDomain.NewService(mockUserRepo)
			ctx := context.Background()
			tt.setup(mockUserRepo)

			err := service.VerifyEmailUniqueness(ctx, tt.email)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnsureUserExists(t *testing.T) {
	var arg = gomock.Any()

	tests := []struct {
		caseName string
		id       idVO.UserID
		setup    func(mockUserRepo *mock.MockIUserRepository)
		errMsg   string
	}{
		{
			caseName: "Positive: ユーザーが存在する場合はエラーが返らない",
			id:       idVO.NewUserIDForTest("existing-user-id"),
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByID(arg, arg).Return(true, nil)
			},
			errMsg: "",
		},
		{
			caseName: "Negative: ユーザーが存在しない場合はエラーが返る",
			id:       idVO.NewUserIDForTest("non-existing-user-id"),
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByID(arg, arg).Return(false, nil)
			},
			errMsg: "user not found",
		},
		{
			caseName: "Negative: ExistsByIDでエラーが返る場合はエラーが返る",
			id:       idVO.NewUserIDForTest("unknown-error-user-id"),
			setup: func(mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().ExistsByID(arg, arg).Return(false, assert.AnError)
			},
			errMsg: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock.NewMockIUserRepository(ctrl)

			service := userDomain.NewService(mockUserRepo)
			ctx := context.Background()
			tt.setup(mockUserRepo)

			err := service.EnsureUserExists(ctx, tt.id)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
