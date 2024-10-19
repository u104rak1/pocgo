package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

func TestEncode(t *testing.T) {
	t.Run("Valid: Hashing good password.", func(t *testing.T) {
		_, err := password.Encode("ValidPassword123")
		assert.NoError(t, err)
	})
	// Since it is difficult to force an error in the GenerateFromPassword function, we omitted testing for errors.
}

func TestCompare(t *testing.T) {
	passwordHash, _ := password.Encode("ValidPassword123")
	tests := []struct {
		caseName string
		password string
		hash     string
		wantErr  error
	}{
		{
			caseName: "Valid: return nil, if the password is matched.",
			password: "ValidPassword123",
			hash:     passwordHash,
			wantErr:  nil,
		},
		{
			caseName: "Invalid: return error, if the password is not matched.",
			password: "DifferentPassword456",
			hash:     passwordHash,
			wantErr:  bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := password.Compare(tt.hash, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
