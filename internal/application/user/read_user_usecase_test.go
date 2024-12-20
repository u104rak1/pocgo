package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userApp "github.com/u104rak1/pocgo/internal/application/user"
	"github.com/u104rak1/pocgo/internal/domain/mock"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestReadUserUsecase(t *testing.T) {
	var (
		userID = idVO.NewUserIDForTest("user").String()
		name   = "Sato Taro"
		email  = "sato@example.com"
		arg    = gomock.Any()
	)

	cmd := userApp.ReadUserCommand{
		ID: userID,
	}

	tests := []struct {
		caseName string
		cmd      userApp.ReadUserCommand
		prepare  func(mockUserServ *mock.MockIUserService, user *userDomain.User)
		wantErr  bool
	}{
		{
			caseName: "Positive: ユーザー取得に成功する",
			cmd:      cmd,
			prepare: func(mockUserServ *mock.MockIUserService, user *userDomain.User) {
				mockUserServ.EXPECT().FindUser(arg, arg).Return(user, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Negative: ユーザーIDが不正な形式である	",
			cmd: userApp.ReadUserCommand{
				ID: "invalid",
			},
			prepare: func(mockUserServ *mock.MockIUserService, user *userDomain.User) {},
			wantErr: true,
		},
		{
			caseName: "Negative: ユーザー取得に失敗する",
			cmd:      cmd,
			prepare: func(mockUserServ *mock.MockIUserService, user *userDomain.User) {
				mockUserServ.EXPECT().FindUser(arg, arg).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserServ := mock.NewMockIUserService(ctrl)
			uc := userApp.NewReadUserUsecase(mockUserServ)
			ctx := context.Background()
			user, err := userDomain.Reconstruct(userID, name, email)
			assert.NoError(t, err)
			tt.prepare(mockUserServ, user)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.Equal(t, userID, dto.ID)
				assert.Equal(t, name, dto.Name)
				assert.Equal(t, email, dto.Email)
			}
		})
	}
}
