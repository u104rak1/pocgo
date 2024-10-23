package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userUC "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/domain/mock"
)

func TestCreateUserUsecase(t *testing.T) {
	validCmd := userUC.CreateUserCommand{
		Name:     "Sato taro",
		Email:    "sato@example.com",
		Password: "password",
	}
	err := errors.New("error")

	type Mocks struct {
		mockUserRepo *mock.MockIUserRepository
		mockAuthRepo *mock.MockIAuthenticationRepository
		mockUserServ *mock.MockIUserService
		mockAuthServ *mock.MockIAuthenticationService
	}

	tests := []struct {
		caseName string
		cmd      userUC.CreateUserCommand
		prepare  func(ctx context.Context, mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "User is successfully created.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			caseName: "Error occurs during VerifyEmailUniqueness in userService.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during userDomain creation.",
			cmd: userUC.CreateUserCommand{
				Email: "sato@example.com",
			},
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during Save in userRepository.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during VerifyUniqueness in authenticationService.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during authenticationDomain creation.",
			cmd: userUC.CreateUserCommand{
				Name:  "Sato taro",
				Email: "sato@example.com",
			},
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during Save in authenticationRepository.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthRepo.EXPECT().Save(ctx, gomock.Any()).Return(err)
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
			mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
			mockUserServ := mock.NewMockIUserService(ctrl)
			mockAuthServ := mock.NewMockIAuthenticationService(ctrl)
			uc := userUC.NewCreateUserUsecase(mockUserRepo, mockAuthRepo, mockUserServ, mockAuthServ)
			ctx := context.Background()
			mocks := Mocks{
				mockUserRepo: mockUserRepo,
				mockAuthRepo: mockAuthRepo,
				mockUserServ: mockUserServ,
				mockAuthServ: mockAuthServ,
			}
			tt.prepare(ctx, mocks)

			dto, err := uc.Run(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dto)
				assert.NotEmpty(t, dto.ID)
				assert.Equal(t, tt.cmd.Name, dto.Name)
				assert.Equal(t, tt.cmd.Email, dto.Email)
			}
		})
	}
}
