package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	authApp "github.com/u104rak1/pocgo/internal/application/authentication"
	"github.com/u104rak1/pocgo/internal/infrastructure/jwt"
)

func TestGenerateAccessToken(t *testing.T) {
	t.Run("アクセストークンを生成できること", func(t *testing.T) {
		t.Parallel()

		service := jwt.NewService()
		userID := "01H2X5JMIN3P8T68PYHXXVK5XN"
		jwtSecretKey := []byte("validSecretKey")

		token, err := service.GenerateAccessToken(userID, jwtSecretKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestGetUserIDFromAccessToken(t *testing.T) {
	tests := []struct {
		name          string
		setupToken    func(service authApp.IJWTService) string
		jwtSecretKey  []byte
		expectedID    string
		expectedError string
	}{
		{
			name: "有効なトークンからユーザーIDを取得できること",
			setupToken: func(service authApp.IJWTService) string {
				token, err := service.GenerateAccessToken("testUserID", []byte("validSecretKey"))
				assert.NoError(t, err)
				return token
			},
			jwtSecretKey:  []byte("validSecretKey"),
			expectedID:    "testUserID",
			expectedError: "",
		},
		{
			name: "無効なトークンの場合エラーを返すこと",
			setupToken: func(service authApp.IJWTService) string {
				return "invalid.token.format"
			},
			jwtSecretKey:  []byte("validSecretKey"),
			expectedID:    "",
			expectedError: "token is malformed: token contains an invalid number of segments",
		},
		{
			name: "署名が無効な場合エラーを返すこと",
			setupToken: func(service authApp.IJWTService) string {
				token, err := service.GenerateAccessToken("testUserID", []byte("correctKey"))
				assert.NoError(t, err)
				return token
			},
			jwtSecretKey:  []byte("wrongKey"),
			expectedID:    "",
			expectedError: "signature is invalid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := jwt.NewService()
			token := tt.setupToken(service)

			userID, err := service.GetUserIDFromAccessToken(token, tt.jwtSecretKey)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
