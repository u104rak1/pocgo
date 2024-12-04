package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userApp "github.com/u104raki/pocgo/internal/application/user"
	"github.com/u104raki/pocgo/internal/domain/mock"
	userDomain "github.com/u104raki/pocgo/internal/domain/user"
	"github.com/u104raki/pocgo/pkg/ulid"
)

func TestReadUserUsecase(t *testing.T) {
	var (
		userID   = ulid.GenerateStaticULID("user")
		userName = "Sato Taro"
		email    = "sato@example.com"
	)
	user, err := userDomain.New(userID, userName, email)
	assert.NoError(t, err)

	cmd := userApp.ReadUserCommand{
		ID: userID,
	}

	tests := []struct {
		caseName string
		cmd      userApp.ReadUserCommand
		prepare  func(ctx context.Context, mockUserRepo *mock.MockIUserRepository)
		wantErr  bool
	}{
		{
			caseName: "User is successfully read.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during FindByID in userRepository.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mockUserRepo *mock.MockIUserRepository) {
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock.NewMockIUserRepository(ctrl)
			uc := userApp.NewReadUserUsecase(mockUserRepo)
			ctx := context.Background()
			tt.prepare(ctx, mockUserRepo)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.Equal(t, userID, dto.ID)
				assert.Equal(t, userName, dto.Name)
				assert.Equal(t, email, dto.Email)
			}
		})
	}
}
