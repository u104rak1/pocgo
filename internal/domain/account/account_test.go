package account_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/account"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/domain/value_object/money"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
	"github.com/ucho456job/pocgo/pkg/timer"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	var (
		accountID = ulid.GenerateStaticULID("account")
		userID    = ulid.GenerateStaticULID("user")
		name      = "For work"
		password  = "1234"
		amount    = 1000.0
		currency  = "JPY"
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
		wantErr   error
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
			wantErr:   nil,
		},
		{
			caseName:  "Error occurs with invalid ID.",
			id:        "invalid",
			userID:    userID,
			name:      name,
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   account.ErrInvalidID,
		},
		{
			caseName:  "Error occurs with invalid UserID.",
			id:        accountID,
			userID:    "invalid",
			name:      name,
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   userDomain.ErrInvalidID,
		},
		{
			caseName:  "Error occurs with 2-character name.",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", account.NameMinLength-1),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   account.ErrInvalidName,
		},
		{
			caseName:  "Successfully creates account with 3-character name.",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", account.NameMinLength),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   nil,
		},
		{
			caseName:  "Successfully creates account with 20-character name.",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", account.NameMaxLength),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   nil,
		},
		{
			caseName:  "Error occurs with 21-character name.",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", account.NameMaxLength+1),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   account.ErrInvalidName,
		},
		{
			caseName:  "Error occurs with 3-character password.",
			id:        accountID,
			userID:    userID,
			name:      strings.Repeat("a", account.NameMaxLength+1),
			password:  password,
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   account.ErrInvalidName,
		},
		{
			caseName:  "an error occurs when the password is 3 characters.",
			id:        accountID,
			userID:    userID,
			name:      name,
			password:  "123",
			amount:    amount,
			currency:  currency,
			updatedAt: now,
			wantErr:   account.ErrPasswordInvalidLength,
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
			wantErr:   nil,
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
			wantErr:   account.ErrPasswordInvalidLength,
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
			wantErr:   money.ErrNegativeAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, err := account.New(tt.id, tt.userID, tt.name, tt.password, tt.amount, tt.currency, tt.updatedAt)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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
		currency  = "JPY"
		now       = timer.GetFixedDate()
	)
	t.Run("Successfully reconstructs an account.", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(password)
		acc, err := account.Reconstruct(accountID, userID, name, encodedPassword, amount, currency, now)

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
		currency  = "JPY"
		now       = timer.Now()
	)

	tests := []struct {
		caseName string
		newName  string
		wantErr  error
	}{
		{
			caseName: "Successfully changes to a valid name.",
			newName:  "NewName",
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with an invalid name.",
			newName:  "",
			wantErr:  account.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ChangeName(tt.newName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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
		currency  = "JPY"
		now       = timer.Now()
	)

	tests := []struct {
		caseName    string
		newPassword string
		wantErr     error
	}{
		{
			caseName:    "Successfully changes to a valid password.",
			newPassword: "5678",
			wantErr:     nil,
		},
		{
			caseName:    "Error occurs with an invalid password.",
			newPassword: "invalid",
			wantErr:     account.ErrPasswordInvalidLength,
		},
		// Since it is difficult to force errors in the Encode function, we have omitted testing for errors.
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ChangePassword(tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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
		currency  = "JPY"
		now       = timer.Now()
	)

	tests := []struct {
		caseName    string
		newPassword string
		wantErr     error
	}{
		{
			caseName:    "Passwords match without errors.",
			newPassword: password,
			wantErr:     nil,
		},
		{
			caseName:    "Error occurs when passwords do not match.",
			newPassword: "invalid",
			wantErr:     account.ErrUnmatchedPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
			err := acc.ComparePassword(tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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
		currency  = "JPY"
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
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with unsupported currency.",
			amount:   300,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
		{
			caseName: "Error occurs when the balance is insufficient.",
			amount:   1500,
			currency: money.JPY,
			wantErr:  money.ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
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
		currency  = "JPY"
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
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with unsupported currency.",
			amount:   300,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
		{
			caseName: "Error occurs when the currency differs.",
			amount:   300,
			currency: money.USD,
			wantErr:  money.ErrDifferentCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
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
		currency  = "JPY"
		now       = timer.Now()
	)

	t.Run("Successfully changes UpdatedAt to valid time.", func(t *testing.T) {
		acc, _ := account.New(accountID, userID, name, password, amount, currency, now)
		newTime := timer.Now()
		acc.ChangeUpdatedAt(newTime)
		assert.Equal(t, newTime, acc.UpdatedAt())
	})
}
