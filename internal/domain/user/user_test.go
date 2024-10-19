package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	validID    = ulid.GenerateStaticULID("user")
	validName  = "sato taro"
	validEmail = "sato@example.com"
)

func TestNew(t *testing.T) {
	tests := []struct {
		caseName string
		id       string
		name     string
		email    string
		wantErr  error
	}{
		{
			caseName: "Happy path: return user entity, if arguments are valid.",
			id:       validID,
			name:     validName,
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "Edge case: return error, if the ID is invalid.",
			id:       "invalid",
			name:     validName,
			email:    validEmail,
			wantErr:  user.ErrInvalidUserID,
		},
		{
			caseName: "Edge case: return error, if the name is 0 characters.",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength-1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "Happy path: return user entity, if the name is 1 characters.",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "Happy path: return user entity, if the name is 20 characters.",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "Edge case: return error, if the name is 21 characters.",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength+1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "Edge case: return error, if the email is invalid.",
			id:       validID,
			name:     validName,
			email:    "invalid",
			wantErr:  user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, err := user.New(tt.id, tt.name, tt.email)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, u.ID())
				assert.Equal(t, tt.name, u.Name())
				assert.Equal(t, tt.email, u.Email())
			}
		})
	}
}

func TestChangeName(t *testing.T) {
	tests := []struct {
		caseName string
		newName  string
		wantErr  error
	}{
		{
			caseName: "Happy path: can be renamed to a valid name.",
			newName:  "yamada hanako",
			wantErr:  nil,
		},
		{
			caseName: "Edge case: return error, invalid name.",
			newName:  strings.Repeat("a", user.NameMaxLength+1),
			wantErr:  user.ErrInvalidUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := user.New(validID, validName, validEmail)
			err := u.ChangeName(tt.newName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validName, u.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, u.Name())
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	tests := []struct {
		caseName string
		newEmail string
		wantErr  error
	}{
		{
			caseName: "Happy path: can be changed to a valid email.",
			newEmail: "yamada@example.com",
			wantErr:  nil,
		},
		{
			caseName: "Edge case: return error, invalid email.",
			newEmail: "invalid-email",
			wantErr:  user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := user.New(validID, validName, validEmail)
			err := u.ChangeEmail(tt.newEmail)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validEmail, u.Email())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newEmail, u.Email())
			}
		})
	}
}
