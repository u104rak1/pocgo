package user

import (
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

type User struct {
	id    idVO.UserID
	name  string
	email string
}

func New(name, email string) (*User, error) {
	id := idVO.NewUserID()
	return newUser(id, name, email)
}

func Reconstruct(id, name, email string) (*User, error) {
	userID, err := idVO.UserIDFromString(id)
	if err != nil {
		return nil, err
	}
	return newUser(userID, name, email)
}

func newUser(id idVO.UserID, name, email string) (*User, error) {
	if err := validName(name); err != nil {
		return nil, err
	}

	if err := validEmail(email); err != nil {
		return nil, err
	}

	return &User{
		id:    id,
		name:  name,
		email: email,
	}, nil
}

func (u *User) ID() idVO.UserID {
	return u.id
}

func (u *User) IDString() string {
	return u.id.String()
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) ChangeName(newName string) error {
	if err := validName(newName); err != nil {
		return err
	}
	u.name = newName
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	if err := validEmail(newEmail); err != nil {
		return err
	}
	u.email = newEmail
	return nil
}
