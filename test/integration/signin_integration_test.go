package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104raki/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104raki/pocgo/internal/presentation/signin"
	"github.com/u104raki/pocgo/pkg/password"
	"github.com/u104raki/pocgo/pkg/ulid"
	"github.com/uptrace/bun"
)

func TestSignin(t *testing.T) {
	var (
		userID         = ulid.GenerateStaticULID("user")
		maxLenUserName = "Sato Taro12345678901"
		email          = "sato@example.com"
		maxLenPassword = "password123456789012"
	)

	prepareNormal := func(t *testing.T, db *bun.DB) {
		user := &model.User{
			ID:    userID,
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
			caseName: "Happy path (201): Signin successfully",
			requestBody: signin.SigninRequest{
				Email:    email,
				Password: maxLenPassword,
			},
			prepare:  prepareNormal,
			wantCode: http.StatusCreated,
		},
		{
			caseName: "Sad path (401): Authentication fails because email not found",
			requestBody: signin.SigninRequest{
				Email:    "diff@example.com",
				Password: maxLenPassword,
			},
			prepare:  prepareNormal,
			wantCode: http.StatusUnauthorized,
		},
		{
			caseName: "Sad path (401): Authentication fails because authentication not found",
			requestBody: signin.SigninRequest{
				Email:    email,
				Password: maxLenPassword,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID,
					Name:  maxLenUserName,
					Email: email,
				}
				InsertTestData(t, db, user)
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			caseName: "Sad path (401): Authentication fails because password is incorrect",
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
