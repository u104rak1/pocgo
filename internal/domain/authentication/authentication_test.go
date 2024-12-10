package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	var (
		userID   = ulid.GenerateStaticULID("user")
		password = "password"
	)

	tests := []struct {
		name     string
		userID   string
		password string
		wantErr  error
	}{
		{
			name:     "Successfully creates an authentication.",
			userID:   userID,
			password: password,
			wantErr:  nil,
		},
		{
			name:     "Error occurs with invalid userID.",
			userID:   "invalid",
			password: password,
			wantErr:  userDomain.ErrInvalidID,
		},
		{
			name:     "Error occurs with 7-character password.",
			userID:   userID,
			password: "1234567",
			wantErr:  authentication.ErrPasswordInvalidLength,
		},
		{
			name:     "Successfully creates an authentication with 8-character password.",
			userID:   userID,
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "Successfully creates an authentication with 20-character password.",
			userID:   userID,
			password: "12345678901234567890",
			wantErr:  nil,
		},
		{
			name:     "Error occurs with 21-character password.",
			userID:   userID,
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
	var (
		userID   = ulid.GenerateStaticULID("user")
		password = "password"
	)
	t.Run("Successfully reconstructs an authentication.", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(password)
		auth, err := authentication.Reconstruct(userID, encodedPassword)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, userID, auth.UserID())
		assert.Equal(t, encodedPassword, auth.PasswordHash())
	})
}

func TestComparePassword(t *testing.T) {
	var (
		userID   = ulid.GenerateStaticULID("user")
		password = "password"
	)

	tests := []struct {
		name        string
		newPassword string
		wantErr     error
	}{
		{
			name:        "Passwords match without errors.",
			newPassword: password,
			wantErr:     nil,
		},
		{
			name:        "Error occurs when passwords do not match.",
			newPassword: "deffirentPassword",
			wantErr:     authentication.ErrUnmatchedPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, _ := authentication.New(userID, password)
			err := auth.ComparePassword(tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
