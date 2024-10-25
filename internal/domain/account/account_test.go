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
	"golang.org/x/crypto/bcrypt"
)

var (
	validID       = ulid.GenerateStaticULID("account")
	validUserID   = ulid.GenerateStaticULID("user")
	validName     = "For work"
	validPassword = "1234"
	validAmount   = 1000.0
	validCurrency = "JPY"
	validTime     = timer.Now()
)

func TestNew(t *testing.T) {
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
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   nil,
		},
		{
			caseName:  "Error occurs with invalid ID.",
			id:        "invalid",
			userID:    validUserID,
			name:      validName,
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrInvalidID,
		},
		{
			caseName:  "Error occurs with invalid UserID.",
			id:        validID,
			userID:    "invalid",
			name:      validName,
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   userDomain.ErrInvalidUserID,
		},
		{
			caseName:  "Error occurs with 2-character name.",
			id:        validID,
			userID:    validUserID,
			name:      strings.Repeat("a", account.NameMinLength-1),
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrInvalidAccountName,
		},
		{
			caseName:  "Successfully creates account with 3-character name.",
			id:        validID,
			userID:    validUserID,
			name:      strings.Repeat("a", account.NameMinLength),
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   nil,
		},
		{
			caseName:  "Successfully creates account with 20-character name.",
			id:        validID,
			userID:    validUserID,
			name:      strings.Repeat("a", account.NameMaxLength),
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   nil,
		},
		{
			caseName:  "Error occurs with 21-character name.",
			id:        validID,
			userID:    validUserID,
			name:      strings.Repeat("a", account.NameMaxLength+1),
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrInvalidAccountName,
		},
		{
			caseName:  "Error occurs with 3-character password.",
			id:        validID,
			userID:    validUserID,
			name:      strings.Repeat("a", account.NameMaxLength+1),
			password:  validPassword,
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrInvalidAccountName,
		},
		{
			caseName:  "an error occurs when the password is 3 characters.",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  "123",
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrPasswordInvalidLength,
		},
		{
			caseName:  "Successfully creates account with 4-character password.",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  "1234",
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   nil,
		},
		{
			caseName:  "Error occurs with 5-character password.",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  "12345",
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrPasswordInvalidLength,
		},
		// Since it is difficult to force errors in the Encode function, we have omitted testing for errors.

		{
			caseName:  "Error occurs with invalid amount.",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  validPassword,
			amount:    -1,
			currency:  validCurrency,
			updatedAt: validTime,
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
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Run("Successfully reconstructs an account.", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(validPassword)
		acc, err := account.Reconstruct(validID, validUserID, validName, encodedPassword, validAmount, validCurrency, validTime)

		assert.NoError(t, err)
		assert.Equal(t, validID, acc.ID())
		assert.Equal(t, validUserID, acc.UserID())
		assert.Equal(t, validName, acc.Name())
		assert.Equal(t, encodedPassword, acc.PasswordHash())
		assert.Equal(t, validAmount, acc.Balance().Amount())
		assert.Equal(t, validCurrency, acc.Balance().Currency())
		assert.Equal(t, validTime, acc.UpdatedAt())
	})
}

func TestChangeName(t *testing.T) {
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
			wantErr:  account.ErrInvalidAccountName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.ChangeName(tt.newName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validName, acc.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, acc.Name())
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
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
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
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
	tests := []struct {
		caseName    string
		newPassword string
		wantErr     error
	}{
		{
			caseName:    "Passwords match without errors.",
			newPassword: validPassword,
			wantErr:     nil,
		},
		{
			caseName:    "Error occurs when passwords do not match.",
			newPassword: "invalid",
			wantErr:     bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
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
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.Withdraw(tt.amount, tt.currency)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, validAmount-tt.amount, acc.Balance().Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
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
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.Deposit(tt.amount, tt.currency)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, validAmount+tt.amount, acc.Balance().Amount())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestChangeUpdatedAt(t *testing.T) {
	t.Run("Successfully changes UpdatedAt to valid time.", func(t *testing.T) {
		acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
		newTime := timer.Now()
		acc.ChangeUpdatedAt(newTime)
		assert.Equal(t, newTime, acc.UpdatedAt())
	})
}
