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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockIUserRepository(ctrl)
	mockAuthRepo := mock.NewMockIAuthenticationRepository(ctrl)
	mockUserServ := mock.NewMockIUserService(ctrl)
	mockAuthServ := mock.NewMockIAuthenticationService(ctrl)
	uc := userUC.NewCreateUserUsecase(mockUserRepo, mockAuthRepo, mockUserServ, mockAuthServ)

	validCmd := userUC.CreateUserCommand{
		Name:     "Sato taro",
		Email:    "sato@example.com",
		Password: "password",
	}
	err := errors.New("error")

	tests := []struct {
		caseName string
		cmd      userUC.CreateUserCommand
		prepare  func(ctx context.Context)
		wantErr  bool
	}{
		{
			caseName: "OK: creates user successfully.",
			cmd:      validCmd,
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mockAuthRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			caseName: "NG: an error occurs userService.VerifyEmailUniqueness.",
			cmd:      validCmd,
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "NG: an error occurs userDomain.New.",
			cmd: userUC.CreateUserCommand{
				Email: "sato@example.com",
			},
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "NG: an error occurs in userRepository.Save.",
			cmd:      validCmd,
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "NG: an error occurs in authenticationService.VerifyUniqueness.",
			cmd:      validCmd,
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
		{
			caseName: "NG: an error occurs in authenticationDomain.New.",
			cmd: userUC.CreateUserCommand{
				Name:  "Sato taro",
				Email: "sato@example.com",
			},
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			caseName: "NG: an error occurs in authenticationRepository.Save.",
			cmd:      validCmd,
			prepare: func(ctx context.Context) {
				mockUserServ.EXPECT().VerifyEmailUniqueness(ctx, gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
				mockAuthServ.EXPECT().VerifyUniqueness(ctx, gomock.Any()).Return(nil)
				mockAuthRepo.EXPECT().Save(ctx, gomock.Any()).Return(err)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			ctx := context.Background()
			tt.prepare(ctx)
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
