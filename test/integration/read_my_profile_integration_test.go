package integration_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func TestReadMyProfile(t *testing.T) {
	var (
		userID         = idVO.NewUserIDForTest("user")
		maxLenUserName = "Sato Taro"
		email          = "sato@example.com"
	)

	tests := []struct {
		caseName string
		prepare  func(t *testing.T, db *bun.DB)
		wantCode int
	}{
		{
			caseName: "Happy path (200): ユーザー情報取得に成功する",
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID.String(),
					Name:  maxLenUserName,
					Email: email,
				}
				InsertTestData(t, db, user)
			},
			wantCode: http.StatusOK,
		},
		{
			caseName: "Sad path (404): ユーザーが見つからない為、失敗する",
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
			SetAccessToken(t, userID.String(), req)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)

			result := GenerateResultJSON(t, beforeDBData, nil, req, rec, nil)
			result = ReplaceDynamicValue(result, []string{})

			gol.Assert(t, t.Name(), result)
		})
	}
}
