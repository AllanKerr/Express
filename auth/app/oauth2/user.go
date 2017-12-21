package oauth2

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/ory/fosite"
)

type User interface {
	// GetID returns the client ID.
	GetUsername() string

	// GetHashedSecret returns the hashed secret as it is stored in the store.
	GetHashedPassword() []byte

	// Returns the scopes this client is allowed to request.
	GetScopes() fosite.Arguments
}

type DefaultUser struct {
	Username string
	PasswordHash []byte
	Scopes []string
}

func (u *DefaultUser) GetUsername() string {
	return u.Username
}

func (u *DefaultUser) GetHashedPassword() []byte {
	return u.PasswordHash
}

func (u *DefaultUser) GetScopes() fosite.Arguments {
	return u.Scopes
}


func NewUser(username string, password string) (*DefaultUser, error) {

	return &DefaultUser{
		username,
		nil,
		nil,
	}, nil
}

func NewAdmin(username string, password string) (*DefaultUser, error) {

	user, err := NewUser(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (user *DefaultUser) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}
