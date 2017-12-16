package oauth2

import "golang.org/x/crypto/bcrypt"

type Secure interface {
	VerifyPassword(password string) error
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
