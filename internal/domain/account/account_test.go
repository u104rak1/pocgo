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
			caseName:  "Successfully creates an account.",
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
			caseName:  "Error occurs with invalid id.",
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
			caseName:  "Error occurs with invalid user id.",
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
			caseName:  "Error occurs with 2-character name.",
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
			caseName:  "Successfully creates account with 3-character name.",
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
			caseName:  "Successfully creates account with 20-character name.",
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
			caseName:  "Error occurs with 21-character name.",
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
			caseName:  "Error occurs with 3-character password.",
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
			caseName:  "Successfully creates account with 4-character password.",
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
			caseName:  "Error occurs with 5-character password.",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  "12345",
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			errMsg:    "account password must be 4 characters",
		},

		// Since it is difficult to force errors in the Encode function, we have omitted testing for errors.

		{
			caseName:  "Error occurs with invalid amount.",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  password,
			amount:    -1,
			currency:  currency,
			updatedAt: now,
			errMsg:    moneyVO.ErrNegativeAmount.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, err := accountDomain.New(tt.id, tt.userID, tt.name, tt.password, tt.amount, tt.currency, tt.updatedAt)

			if tt.errMsg != "" {
				assert.Equal(t, err.Error(), tt.errMsg)
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
	t.Run("Successfully reconstructs an account.", func(t *testing.T) {
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
			caseName: "Successfully changes to a valid name.",
			newName:  "NewName",
			errMsg:   "",
		},
		{
			caseName: "Error occurs with an invalid name.",
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
				assert.Equal(t, err.Error(), tt.errMsg)
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
			caseName:    "Successfully changes to a valid password.",
			newPassword: "5678",
			errMsg:      "",
		},
		{
			caseName:    "Error occurs with an invalid password.",
			newPassword: "invalid",
			errMsg:      "account password must be 4 characters",
		},
		// Since it is difficult to force errors in the Encode function, we have omitted testing for errors.
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ChangePassword(tt.newPassword)

			if tt.errMsg != "" {
				assert.Equal(t, err.Error(), tt.errMsg)
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
			caseName:    "Passwords match without errors.",
			newPassword: password,
			errMsg:      "",
		},
		{
			caseName:    "Error occurs when passwords do not match.",
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
				assert.Equal(t, err.Error(), tt.errMsg)
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
		wantErr  error
	}{
		{
			caseName: "Successfully withdraws when the currency matches and the balance is sufficient.",
			amount:   300,
			currency: moneyVO.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with unsupported currency.",
			amount:   300,
			currency: "EUR",
			wantErr:  moneyVO.ErrUnsupportedCurrency,
		},
		{
			caseName: "Error occurs when the balance is insufficient.",
			amount:   1500,
			currency: moneyVO.JPY,
			wantErr:  moneyVO.ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.Withdraw(tt.amount, tt.currency)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, amount-tt.amount, acc.Balance().Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
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
		wantErr  error
	}{
		{
			caseName: "Successfully deposits when the currency matches.",
			amount:   300,
			currency: moneyVO.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with unsupported currency.",
			amount:   300,
			currency: "EUR",
			wantErr:  moneyVO.ErrUnsupportedCurrency,
		},
		{
			caseName: "Error occurs when the currency differs.",
			amount:   300,
			currency: moneyVO.USD,
			wantErr:  moneyVO.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
			err := acc.Deposit(tt.amount, tt.currency)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, amount+tt.amount, acc.Balance().Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
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

	t.Run("Successfully changes UpdatedAt to valid time.", func(t *testing.T) {
		acc, _ := accountDomain.New(accountID, userID, name, password, amount, currency, now)
		newTime := timer.Now()
		acc.ChangeUpdatedAt(newTime)
		assert.Equal(t, newTime, acc.UpdatedAt())
	})
}
