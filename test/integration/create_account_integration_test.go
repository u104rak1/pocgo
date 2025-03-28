package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/domain/value_object/money"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/seed"
	"github.com/u104rak1/pocgo/internal/presentation/me/accounts"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/uptrace/bun"
)

func TestCreateAccount(t *testing.T) {
	var (
		name      = "AccountName123456789"
		password  = "1234"
		currency  = money.JPY
		userID    = idVO.NewUserIDForTest("user")
		userName  = "sato taro"
		userEmail = "sata@example.com"
	)

	tests := []struct {
		caseName    string
		requestBody interface{}
		prepare     func(t *testing.T, db *bun.DB)
		wantCode    int
	}{
		{
			caseName: "Happy path (201): 口座作成に成功する",
			requestBody: accounts.CreateAccountRequestBody{
				Name:     name,
				Password: password,
				Currency: currency,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID.String(),
					Name:  userName,
					Email: userEmail,
				}
				InsertTestData(t, db, user)
			},
			wantCode: http.StatusCreated,
		},
		{
			caseName: "Sad path (404): ユーザーが見つからない為、失敗する",
			requestBody: accounts.CreateAccountRequestBody{
				Name:     name,
				Password: password,
				Currency: currency,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				InsertTestData(t, db)
			},
			wantCode: http.StatusNotFound,
		},
		{
			caseName: "Sad path (409): 口座の上限に達している為、失敗する",
			requestBody: accounts.CreateAccountRequestBody{
				Name:     name,
				Password: password,
				Currency: currency,
			},
			prepare: func(t *testing.T, db *bun.DB) {
				user := &model.User{
					ID:    userID.String(),
					Name:  userName,
					Email: userEmail,
				}
				var accounts []*model.Account
				for i := 0; i < accountDomain.MaxAccountLimit; i++ {
					passwordHash, err := passwordUtil.Encode("1234")
					assert.NoError(t, err)
					accounts = append(accounts, &model.Account{
						ID:           idVO.NewAccountIDForTest(fmt.Sprintf("account%d", i)).String(),
						UserID:       userID.String(),
						Name:         fmt.Sprintf("Account%d", i),
						PasswordHash: passwordHash,
						Balance:      0,
						CurrencyID:   seed.JPYID,
						UpdatedAt:    timer.Now(),
					})
				}
				InsertTestData(t, db, user, accounts)
			},
			wantCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			e, gol, db := BeforeAll(t)
			defer AfterAll(t, db)

			tt.prepare(t, db)
			usedTables := []string{"users", "accounts"}
			beforeDBData := GetDBData(t, db, usedTables)

			req, rec := NewJSONRequest(t, http.MethodPost, "/api/v1/me/accounts", tt.requestBody)
			SetAccessToken(t, userID.String(), req)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantCode, rec.Code)

			afterDBData := GetDBData(t, db, usedTables)
			result := GenerateResultJSON(t, beforeDBData, afterDBData, req, rec, tt.requestBody)
			replaceKeys := []string{"id", "passwordHash", "accessToken", "updatedAt"}
			result = ReplaceDynamicValue(result, replaceKeys)

			gol.Assert(t, t.Name(), result)
		})
	}
}
