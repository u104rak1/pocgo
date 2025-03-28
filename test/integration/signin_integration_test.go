package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/internal/presentation/signin"
	"github.com/u104rak1/pocgo/pkg/password"
	"github.com/uptrace/bun"
)

func TestSignin(t *testing.T) {
	var (
		userID         = idVO.NewUserIDForTest("user")
		maxLenUserName = "Sato Taro12345678901"
		email          = "sato@example.com"
		maxLenPassword = "password123456789012"
	)

	prepareNormal := func(t *testing.T, db *bun.DB) {
		user := &model.User{
			ID:    userID.String(),
			Name:  maxLenUserName,
			Email: email,
		}
		passwordHash, err := password.Encode(maxLenPassword)
		assert.NoError(t, err)
		auth := &model.Authentication{
			UserID:       user.ID,
			PasswordHash: passwordHash,
		}
		InsertTestData(t, db, user, auth)
	}

	tests := []struct {
		caseName    string
		requestBody interface{}
		prepare     func(t *testing.T, db *bun.DB)
		wantCode    int
	}{
		{
			caseName: "Happy path (201): ログインに成功する",
			requestBody: signin.SigninRequest{
				Email:    email,
				Password: maxLenPassword,
			},
			prepare:  prepareNormal,
			wantCode: http.StatusCreated,
		},
		{
			caseName: "Sad path (401): 指定したメールアドレスを持つユーザーが存在しない為、失敗する",
			requestBody: signin.SigninRequest{
				Email:    "diff@example.com",
				Password: maxLenPassword,
			},
			prepare:  prepareNormal,
			wantCode: http.StatusUnauthorized,
		},
		{
			caseName: "Sad path (401): 認証情報が見つからない為、失敗する",
			requestBody: signin.SigninRequest{
				Email:    email,
				Password: maxLenPassword,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID.String(),
					Name:  maxLenUserName,
					Email: email,
				}
				InsertTestData(t, db, user)
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			caseName: "Sad path (401): パスワードが異なる為、失敗する",
			requestBody: signin.SigninRequest{
				Email:    email,
				Password: "diffPassword",
			},
			prepare:  prepareNormal,
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			e, gol, db := BeforeAll(t)
			defer AfterAll(t, db)

			tt.prepare(t, db)
			usedTables := []string{"users", "authentications"}
			beforeDBData := GetDBData(t, db, usedTables)

			req, rec := NewJSONRequest(t, http.MethodPost, "/api/v1/signin", tt.requestBody)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)

			result := GenerateResultJSON(t, beforeDBData, nil, req, rec, tt.requestBody)
			replaceKeys := []string{"passwordHash", "accessToken"}
			result = ReplaceDynamicValue(result, replaceKeys)

			gol.Assert(t, t.Name(), result)
		})
	}
}
