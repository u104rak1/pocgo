package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/internal/presentation/signup"
	"github.com/u104rak1/pocgo/pkg/ulid"
	"github.com/uptrace/bun"
)

func TestSignup(t *testing.T) {
	var (
		maxLenUserName     = "Sato Taro12345678901"
		userEmail          = "sato@example.com"
		maxLenUserPassword = "password123456789012"
	)

	tests := []struct {
		caseName    string
		requestBody interface{}
		prepare     func(t *testing.T, db *bun.DB)
		wantCode    int
	}{
		{
			caseName: "Happy path (201): Signup successfully",
			requestBody: signup.SignupRequest{
				Name:     maxLenUserName,
				Email:    userEmail,
				Password: maxLenUserPassword,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				InsertTestData(t, db)
			},
			wantCode: http.StatusCreated,
		},
		{
			caseName: "Sad path (409): email is already used",
			requestBody: signup.SignupRequest{
				Name:     maxLenUserName,
				Email:    "conflict@example.com",
				Password: maxLenUserPassword,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				existingUser := &model.User{
					ID:    ulid.GenerateStaticULID("user"),
					Name:  "Existing User",
					Email: "conflict@example.com",
				}
				InsertTestData(t, db, existingUser)
			},
			wantCode: http.StatusConflict,
		},
		// Exclude duplicate error of authentication because they occur infrequently and are difficult to reproduce.
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			e, gol, db := BeforeAll(t)
			defer AfterAll(t, db)

			tt.prepare(t, db)
			usedTables := []string{"users", "authentications"}
			beforeDBData := GetDBData(t, db, usedTables)

			req, rec := NewJSONRequest(t, http.MethodPost, "/api/v1/signup", tt.requestBody)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)

			afterDBData := GetDBData(t, db, usedTables)
			result := GenerateResultJSON(t, beforeDBData, afterDBData, req, rec, tt.requestBody)
			replaceKeys := []string{"id", "userId", "passwordHash", "updatedAt", "accessToken"}
			result = ReplaceDynamicValue(result, replaceKeys)

			gol.Assert(t, t.Name(), result)
		})
	}
}
