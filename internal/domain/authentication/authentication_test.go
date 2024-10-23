package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
	"github.com/ucho456job/pocgo/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

var (
	validUserID   = ulid.GenerateStaticULID("user")
	validPassword = "password"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		password string
		wantErr  error
	}{
		{
			name:     "Successfully creates an authentication.",
			userID:   validUserID,
			password: validPassword,
			wantErr:  nil,
		},
		{
			name:     "Error occurs with invalid userID.",
			userID:   "invalid",
			password: validPassword,
			wantErr:  userDomain.ErrInvalidUserID,
		},
		{
			name:     "Error occurs with 7-character password.",
			userID:   validUserID,
			password: "1234567",
			wantErr:  authentication.ErrPasswordInvalidLength,
		},
		{
			name:     "Successfully creates an authentication with 8-character password.",
			userID:   validUserID,
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "Successfully creates an authentication with 20-character password.",
			userID:   validUserID,
			password: "12345678901234567890",
			wantErr:  nil,
		},
		{
			name:     "Error occurs with 21-character password.",
			userID:   validUserID,
			password: "123456789012345678901",
			wantErr:  authentication.ErrPasswordInvalidLength,
		},
		// Since it is difficult to force errors in the Encode function, we have omitted testing for errors.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, err := authentication.New(tt.userID, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, auth)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, auth.UserID())
				assert.NoError(t, passwordUtil.Compare(auth.PasswordHash(), tt.password))
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Run("Successfully reconstructs an authentication.", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(validPassword)
		auth, err := authentication.Reconstruct(validUserID, encodedPassword)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, validUserID, auth.UserID())
		assert.Equal(t, encodedPassword, auth.PasswordHash())
	})
}

func TestComparePassword(t *testing.T) {
	tests := []struct {
		name        string
		newPassword string
		wantErr     error
	}{
		{
			name:        "Passwords match without errors.",
			newPassword: validPassword,
			wantErr:     nil,
		},
		{
			name:        "Error occurs when passwords do not match.",
			newPassword: "deffirentPassword",
			wantErr:     bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, _ := authentication.New(validUserID, validPassword)
			err := auth.ComparePassword(tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
