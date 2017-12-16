package oauth2

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	PasswordHash []byte
	Role string
}

func NewUser(username string, password string) (*User, error) {

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{username, passwordHash, "user"}, nil
}

func NewAdmin(username string, password string) (*User, error) {

	user, err := NewUser(username, password)
	if err != nil {
		return nil, err
	}
	user.Role = "admin"
	return user, nil
}

func (user *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}
