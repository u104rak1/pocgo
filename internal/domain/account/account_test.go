package account_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/timer"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.GetFixedDate()
	)

	tests := []struct {
		caseName  string
		id        string
		userID    string
		name      string
		password  string
		amount    float64
		currency  string
		updatedAt time.Time
		errMsg    string
	}{
		{
			caseName:  "Positive: 口座を作成できる",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "",
		},
		{
			caseName:  "Negative: 無効なIDの場合はエラーが返る",
			id:        "invalid",
			userID:    userID,
			name:      name,
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account id must be a valid ULID",
		},
		{
			caseName:  "Negative: 無効なユーザーIDの場合はエラーが返る",
			id:        accountID,
			userID:    "invalid",
			name:      name,
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    userDomain.ErrInvalidID.Error(),
		},
		{
			caseName:  "Negative: 2文字の名前の場合はエラーが返る",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", 3-1),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account name must be between 3 and 20 characters",
		},
		{
			caseName:  "Positive: 3文字の名前の場合は口座を作成できる",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", 3),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "",
		},
		{
			caseName:  "Positive: 20文字の名前の場合は口座を作成できる",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", 20),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "",
		},
		{
			caseName:  "Negative: 21文字の名前の場合はエラーが返る",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", 20+1),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account name must be between 3 and 20 characters",
		},
		{
			caseName:  "Negative: 3文字のパスワードの場合はエラーが返る",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  "123",
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account password must be 4 characters",
		},
		{
			caseName:  "Positive: 4文字のパスワードの場合は口座を作成できる",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  "1234",
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "",
		},
		{
			caseName:  "Negative: 5文字のパスワードの場合はエラーが返る",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  "12345",
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account password must be 4 characters",
		},
		{
			caseName:  "Negative: 無効な金額の場合はエラーが返る",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  password,
			amount:    -1,
			currency:  currency,
			updatedAt: now,
			errMsg:    moneyVO.ErrNegativeAmount.Error(),
		},

		// Password.Encode関数を強制的にエラーにすることが難しい為、このエラーパターンはテストしない
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, err := accountDomain.New(tt.id, tt.userID, tt.name, tt.password, tt.amount, tt.currency, tt.updatedAt)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, acc)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, acc.ID())
				assert.Equal(t, tt.userID, acc.UserID())
				assert.Equal(t, tt.name, acc.Name())
				assert.NoError(t, passwordUtil.Compare(acc.PasswordHash(), tt.password))
				assert.Equal(t, tt.amount, acc.Balance().Amount())
				assert.Equal(t, tt.currency, acc.Balance().Currency())
				assert.Equal(t, tt.updatedAt, acc.UpdatedAt())
				assert.Equal(t, timer.GetFixedDateString(), acc.UpdatedAtString())
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.GetFixedDate()
	)
	t.Run("Positive: 口座を再構築できる", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(password)
		acc, err := accountDomain.Reconstruct(accountID, userID, name, encodedPassword, amount, currency, now)

		assert.NoError(t, err)
		assert.Equal(t, accountID, acc.ID())
		assert.Equal(t, userID, acc.UserID())
		assert.Equal(t, name, acc.Name())
		assert.Equal(t, encodedPassword, acc.PasswordHash())
		assert.Equal(t, amount, acc.Balance().Amount())
		assert.Equal(t, currency, acc.Balance().Currency())
		assert.Equal(t, now, acc.UpdatedAt())
		assert.Equal(t, timer.GetFixedDateString(), acc.UpdatedAtString())
	})
}

func TestChangeName(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	tests := []struct {
		caseName string
		newName  string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効な名前の場合は名前を変更できる",
			newName:  "NewName",
			errMsg:   "",
		},
		{
			caseName: "Negative: 無効な名前の場合はエラーが返る",
			newName:  "",
			errMsg:   "account name must be between 3 and 20 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ChangeName(tt.newName)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, name, acc.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, acc.Name())
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	tests := []struct {
		caseName    string
		newPassword string
		errMsg      string
	}{
		{
			caseName:    "Positive: 有効なパスワードの場合はパスワードを変更できる",
			newPassword: "5678",
			errMsg:      "",
		},
		{
			caseName:    "Negative: 無効なパスワードの場合はエラーが返る",
			newPassword: "invalid",
			errMsg:      "account password must be 4 characters",
		},
		// Password.Encode関数を強制的にエラーにすることが難しい為、このエラーパターンはテストしない
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ChangePassword(tt.newPassword)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Error(t, passwordUtil.Compare(acc.PasswordHash(), tt.newPassword))
			} else {
				assert.NoError(t, err)
				assert.NoError(t, passwordUtil.Compare(acc.PasswordHash(), tt.newPassword))
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	tests := []struct {
		caseName    string
		newPassword string
		errMsg      string
	}{
		{
			caseName:    "Positive: パスワードが一致している場合はエラーが返らない",
			newPassword: password,
			errMsg:      "",
		},
		{
			caseName:    "Negative: パスワードが一致していない場合はエラーが返る",
			newPassword: "invalid",
			errMsg:      "passwords do not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ComparePassword(tt.newPassword)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	tests := []struct {
		caseName string
		amount   float64
		currency string
		errMsg   string
	}{
		{
			caseName: "Positive: 通貨が一致し、残高が十分な場合は引き出しができる",
			amount:   300,
			currency: moneyVO.JPY,
			errMsg:   "",
		},
		{
			caseName: "Negative: money値オブジェクトの作成に失敗した場合、エラーが返る",
			amount:   300,
			currency: "EUR",
			errMsg:   moneyVO.ErrUnsupportedCurrency.Error(),
		},
		{
			caseName: "Negative: money値オブジェクトのSubメソッドが失敗した場合、エラーが返る",
			amount:   1500,
			currency: moneyVO.JPY,
			errMsg:   moneyVO.ErrInsufficientBalance.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.Withdraw(tt.amount, tt.currency)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, amount-tt.amount, acc.Balance().Amount())
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	tests := []struct {
		caseName string
		amount   float64
		currency string
		errMsg   string
	}{
		{
			caseName: "Positive: 通貨が一致している場合は入金ができる",
			amount:   300,
			currency: moneyVO.JPY,
			errMsg:   "",
		},
		{
			caseName: "Negative: money値オブジェクトの作成に失敗した場合、エラーが返る",
			amount:   300,
			currency: "EUR",
			errMsg:   moneyVO.ErrUnsupportedCurrency.Error(),
		},
		{
			caseName: "Negative: money値オブジェクトのAddメソッドが失敗した場合、エラーが返る",
			amount:   300,
			currency: moneyVO.USD,
			errMsg:   moneyVO.ErrAddDifferentCurrency.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.Deposit(tt.amount, tt.currency)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, amount+tt.amount, acc.Balance().Amount())
			}
		})
	}
}

func TestChangeUpdatedAt(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = moneyVO.JPY
		now       = timer.Now()
	)

	t.Run("Positive: UpdatedAtを有効な時間に変更できる", func(t *testing.T) {
		acc, err := accountDomain.New(accountID, userID, name, password, amount, currency, now)
		assert.NoError(t, err)
		newTime := timer.Now()
		acc.ChangeUpdatedAt(newTime)
		assert.Equal(t, acc.UpdatedAt(), newTime)
	})
}
