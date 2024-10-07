package user_domain

type User struct {
	id    string
	name  string
	email string
}

func New(id, name, email string) (*User, error) {
	var err error
	if err = IsValidID(id); err != nil {
		return nil, err
	}

	if err = isValidName(name); err != nil {
		return nil, err
	}

	if err = isValidEmail(email); err != nil {
		return nil, err
	}

	return &User{
		id:    id,
		name:  name,
		email: email,
	}, nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) ChangeName(newName string) error {
	if err := isValidName(newName); err != nil {
		return err
	}
	u.name = newName
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	if err := isValidEmail(newEmail); err != nil {
		return err
	}
	u.email = newEmail
	return nil
}
