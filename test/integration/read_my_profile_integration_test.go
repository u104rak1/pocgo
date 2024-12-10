package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/pkg/ulid"
	"github.com/uptrace/bun"
)

func TestReadMyProfileHandler(t *testing.T) {
	var (
		userID         = ulid.GenerateStaticULID("user")
		maxLenUserName = "Sato Taro"
		email          = "sato@example.com"
	)

	tests := []struct {
		caseName string
		prepare  func(t *testing.T, db *bun.DB)
		wantCode int
	}{
		{
			caseName: "Happy path (200): Successful profile retrieval",
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID,
					Name:  maxLenUserName,
					Email: email,
				}
				InsertTestData(t, db, user)
			},
			wantCode: http.StatusOK,
		},
		{
			caseName: "Sad path (404): User not found",
			prepare: func(t *testing.T, db *bun.DB) {
				InsertTestData(t, db)
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			e, gol, db := BeforeAll(t)
			defer AfterAll(t, db)

			tt.prepare(t, db)
			usedTables := []string{"users"}
			beforeDBData := GetDBData(t, db, usedTables)

			req, rec := NewJSONRequest(t, http.MethodGet, "/api/v1/me", nil)
			SetAccessToken(t, userID, req)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)

			result := GenerateResultJSON(t, beforeDBData, nil, req, rec, nil)
			result = ReplaceDynamicValue(result, []string{})

			gol.Assert(t, t.Name(), result)
		})
	}
}
