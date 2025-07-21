package app

import (
	"golang.org/x/crypto/bcrypt"
)

// Interface PasswordHasher agar mudah diuji & swap ke algoritma lain
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

type BcryptHasher struct{}

func (b *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
