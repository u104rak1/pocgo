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
			caseName: "Positive: creates user successfully.",
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
			caseName: "Negative: an error occurs in userService.VerifyEmailUniqueness.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: an error occurs in userDomain.New.",
			cmd: userUC.CreateUserCommand{
				Email: "sato@example.com",
			},
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: an error occurs in userRepository.Save.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: an error occurs in authenticationService.VerifyUniqueness.",
			cmd:      validCmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "Negative: an error occurs in authenticationDomain.New.",
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
			caseName: "Negative: an error occurs in authenticationRepository.Save.",
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