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
	"github.com/ucho456job/pocgo/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

var (
	validID       = ulid.GenerateStaticULID("valid")
	validUserID   = ulid.GenerateStaticULID("user")
	validName     = "For work"
	validPassword = "1234"
	validAmount   = 1000.0
	validCurrency = "JPY"
	validTime     = time.Now()
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
			caseName:  "Happy path: 有効なAccountエンティティを作成",
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
			caseName:  "Edge case: 無効なIDを指定するとエラー",
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
			caseName:  "Edge case: 無効なユーザーIDを指定するとエラー",
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
			caseName:  "Edge case: 名前が0文字だとエラー",
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
			caseName:  "Happy path: 名前が1文字なら成功",
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
			caseName:  "Happy path: 名前が10文字なら成功",
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
			caseName:  "Happy path: 名前が11文字だとエラー",
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
			caseName:  "Edge case: パスワードが3文字だとエラー",
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
			caseName:  "Happy path: パスワードが4文字なら成功",
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
			caseName:  "Edge case: パスワードが5文字だとエラー",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  "12345",
			amount:    validAmount,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   account.ErrPasswordInvalidLength,
		},
		{
			caseName:  "Edge case: 負の残高を指定するとエラー",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  validPassword,
			amount:    -1,
			currency:  validCurrency,
			updatedAt: validTime,
			wantErr:   money.ErrNegativeAmount,
		},
		{
			caseName:  "Edge case: 未対応の通貨を指定するとエラー",
			id:        validID,
			userID:    validUserID,
			name:      validName,
			password:  validPassword,
			amount:    validAmount,
			currency:  "EUR",
			updatedAt: validTime,
			wantErr:   money.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, err := account.New(tt.id, tt.userID, tt.name, tt.password, tt.amount, tt.currency, tt.updatedAt)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, acc.ID())
				assert.Equal(t, tt.userID, acc.UserID())
				assert.Equal(t, tt.name, acc.Name())
				assert.NoError(t, passwordUtil.Compare(acc.PasswordHash(), tt.password))
				assert.Equal(t, tt.amount, acc.Balance().Amount())
				assert.Equal(t, tt.currency, acc.Balance().Currency())
				assert.Equal(t, tt.updatedAt, acc.UpdatedAt())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, acc)
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Run("Happy path: 有効なAccountエンティティを再構築", func(t *testing.T) {
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
			caseName: "Happy path: 名前を変更",
			newName:  "NewName",
			wantErr:  nil,
		},
		{
			caseName: "Edge case: 無効な名前を指定",
			newName:  "",
			wantErr:  account.ErrInvalidAccountName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.ChangeName(tt.newName)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, acc.Name())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validName, acc.Name())
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
			caseName:    "Happy path: 有効なパスワードに変更できる",
			newPassword: "5678",
			wantErr:     nil,
		},
		{
			caseName:    "Edge case: 無効なパスワードを指定した時、パスワードを変更できない",
			newPassword: "invalid",
			wantErr:     account.ErrPasswordInvalidLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.ChangePassword(tt.newPassword)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.NoError(t, passwordUtil.Compare(acc.PasswordHash(), tt.newPassword))
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Error(t, passwordUtil.Compare(acc.PasswordHash(), tt.newPassword))
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
			caseName:    "Happy path: パスワードが一致する場合、エラーが発生しない",
			newPassword: "1234",
			wantErr:     nil,
		},
		{
			caseName:    "Edge case: パスワードが異なる場合、エラーが発生する",
			newPassword: "invalid",
			wantErr:     bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
			err := acc.ComparePassword(tt.newPassword)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
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
			caseName: "Happy path: 通貨が一致し残高が十分ある場合、残高が引かれる",
			amount:   300,
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Edge case: 未対応の通貨の場合、エラーが発生する",
			amount:   300,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
		{
			caseName: "Edge case: 残高が不十分な場合、エラーが発生する",
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
			caseName: "Happy path: 通貨が一致する場合、残高に入金できる",
			amount:   300,
			currency: money.JPY,
			wantErr:  nil,
		},
		{
			caseName: "Edge case: 未対応の通貨の場合、エラーが発生する",
			amount:   300,
			currency: "EUR",
			wantErr:  money.ErrUnsupportedCurrency,
		},
		{
			caseName: "Edge case: 通貨が異なる場合、エラーが発生する",
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
	t.Run("Happy path: 有効な更新日時に変更できる", func(t *testing.T) {
		acc, _ := account.New(validID, validUserID, validName, validPassword, validAmount, validCurrency, validTime)
		newTime := time.Now()
		acc.ChangeUpdatedAt(newTime)
		assert.Equal(t, newTime, acc.UpdatedAt())
	})
}
