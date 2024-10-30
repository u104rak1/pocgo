package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	var (
		id    = ulid.GenerateStaticULID("user")
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		id       string
		name     string
		email    string
		wantErr  error
	}{
		{
			caseName: "Successfully creates a user.",
			id:       id,
			name:     name,
			email:    email,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with invalid ID.",
			id:       "invalid",
			name:     name,
			email:    email,
			wantErr:  user.ErrInvalidID,
		},
		{
			caseName: "Error occurs with 2-character name.",
			id:       id,
			name:     strings.Repeat("a", user.NameMinLength-1),
			email:    email,
			wantErr:  user.ErrInvalidName,
		},
		{
			caseName: "Successfully creates a user with 3-character name.",
			id:       id,
			name:     strings.Repeat("a", user.NameMinLength),
			email:    email,
			wantErr:  nil,
		},
		{
			caseName: "Successfully creates a user with 20-character name.",
			id:       id,
			name:     strings.Repeat("a", user.NameMaxLength),
			email:    email,
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with 21-character name.",
			id:       id,
			name:     strings.Repeat("a", user.NameMaxLength+1),
			email:    email,
			wantErr:  user.ErrInvalidName,
		},
		{
			caseName: "Error occurs with invalid email.",
			id:       id,
			name:     name,
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
	var (
		id    = ulid.GenerateStaticULID("user")
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		newName  string
		wantErr  error
	}{
		{
			caseName: "Successfully changes to a valid name.",
			newName:  "yamada hanako",
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with an invalid name.",
			newName:  strings.Repeat("a", user.NameMaxLength+1),
			wantErr:  user.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := user.New(id, name, email)
			err := u.ChangeName(tt.newName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, name, u.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, u.Name())
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	var (
		id    = ulid.GenerateStaticULID("user")
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		newEmail string
		wantErr  error
	}{
		{
			caseName: "Successfully changes to a valid email.",
			newEmail: "yamada@example.com",
			wantErr:  nil,
		},
		{
			caseName: "Error occurs with an invalid email.",
			newEmail: "invalid-email",
			wantErr:  user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := user.New(id, name, email)
			err := u.ChangeEmail(tt.newEmail)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, email, u.Email())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newEmail, u.Email())
			}
		})
	}
}
