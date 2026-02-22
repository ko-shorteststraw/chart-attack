package entities

type ValidatedUser struct {
	user *User
}

func NewValidatedUser(u *User) (*ValidatedUser, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}
	return &ValidatedUser{user: u}, nil
}

func (v *ValidatedUser) User() *User {
	return v.user
}
