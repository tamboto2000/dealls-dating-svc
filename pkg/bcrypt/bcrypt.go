package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = 10
)

var ErrHashNotMatch = errors.New("hash is not a match with provided raw data")

func Hash(b []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(b, cost)
}

func HashCompare(hashed []byte, b []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashed, b); err != nil {
		return ErrHashNotMatch
	}

	return nil
}
