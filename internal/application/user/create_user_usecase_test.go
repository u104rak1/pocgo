package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userApp "github.com/ucho456job/pocgo/internal/application/user"
	"github.com/ucho456job/pocgo/internal/domain/mock"
)

func TestCreateUserUsecase(t *testing.T) {
	type Mocks struct {
		mockUserRepo *mock.MockIUserRepository
		mockAuthRepo *mock.MockIAuthenticationRepository
		mockUserServ *mock.MockIUserService
		mockAuthServ *mock.MockIAuthenticationService
	}

	var (
		name     = "Sato taro"
		email    = "sato@example.com"
		password = "password"
	)

	cmd := userApp.CreateUserCommand{
		Name:     name,
		Email:    email,
		Password: password,
	}

	tests := []struct {
		caseName string
		cmd      userApp.CreateUserCommand
		prepare  func(ctx context.Context, mocks Mocks)
		wantErr  bool
	}{
		{
			caseName: "User is successfully created.",
			cmd:      cmd,
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
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during userDomain creation.",
			cmd: userApp.CreateUserCommand{
				Email: email,
			},
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during Save in userRepository.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during VerifyUniqueness in authenticationService.",
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			caseName: "Error occurs during authenticationDomain creation.",
			cmd: userApp.CreateUserCommand{
				Name:  name,
				Email: email,
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
			cmd:      cmd,
			prepare: func(ctx context.Context, mocks Mocks) {
				mocks.mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mocks.mockAuthRepo.EXPECT().Save(ctx, gomock.Any()).Return(errors.New("error"))
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
			uc := userApp.NewCreateUserUsecase(mockUserRepo, mockAuthRepo, mockUserServ, mockAuthServ)
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
