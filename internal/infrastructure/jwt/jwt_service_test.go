package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/jwt"
)

func TestGenerateAccessToken(t *testing.T) {
	t.Run("アクセストークンを生成できること", func(t *testing.T) {
		t.Parallel()

		service := jwt.NewService([]byte("validSecretKey"))
		userID := "01H2X5JMIN3P8T68PYHXXVK5XN"

		token, err := service.GenerateAccessToken(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestGetUserIDFromAccessToken(t *testing.T) {
	var (
		userID = idVO.NewUserIDForTest("testUserID")
	)

	tests := []struct {
		name         string
		setupToken   func(service authApp.IJWTService) string
		jwtSecretKey []byte
		expectedID   string
		errMsg       string
	}{
		{
			name: "有効なトークンからユーザーIDを取得できること",
			setupToken: func(service authApp.IJWTService) string {
				token, err := service.GenerateAccessToken(userID.String())
				assert.NoError(t, err)
				return token
			},
			jwtSecretKey: []byte("validSecretKey"),
			expectedID:   userID.String(),
			errMsg:       "",
		},
		{
			name: "無効なトークンの場合エラーを返すこと",
			setupToken: func(service authApp.IJWTService) string {
				return "invalid.token.format"
			},
			jwtSecretKey: []byte("validSecretKey"),
			expectedID:   "",
			errMsg:       "token is malformed:",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := jwt.NewService(tt.jwtSecretKey)
			token := tt.setupToken(service)

			userID, err := service.GetUserIDFromAccessToken(token)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
