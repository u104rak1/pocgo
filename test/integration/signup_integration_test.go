package integration_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/presentation/signup"
	"github.com/uptrace/bun"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		caseName    string
		requestBody signup.SignupRequestBody
		prepare     func(t *testing.T, db *bun.DB)
		wantCode    int
	}{
		{
			caseName: "Happy path: Signup success",
			requestBody: signup.SignupRequestBody{
				User: signup.SignupRequestBodyUser{
					Name:     "Sato Taro",
					Email:    "sato@example.com",
					Password: "password",
					Account: signup.SignupRequestBodyAccount{
						Name:     "For work",
						Password: "1234",
						Currency: "JPY",
					},
				},
			},
			prepare: func(t *testing.T, db *bun.DB) {
				InsertTestData(t, db)
			},
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			e, gol, db := BeforeAll(t)
			defer AfterAll(t, db)

			tt.prepare(t, db)
			usedTables := []string{"users", "accounts", "authentications"}
			beforeDBData := GetDBData(t, db, usedTables)

			req, rec := NewJSONRequest(t, http.MethodPost, "/api/v1/signup", tt.requestBody)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)

			afterDBData := GetDBData(t, db, usedTables)
			result := GenerateResultJSON(t, beforeDBData, afterDBData, req, rec)

			ulidPattern := regexp.MustCompile(`[0-9A-HJKMNP-TV-Z]{26}`)
			modifiedJSON := ulidPattern.ReplaceAll(result, []byte("ANY_ID"))

			gol.Assert(t, t.Name(), modifiedJSON)
		})
	}
}
