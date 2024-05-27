package secrand

import (
	"crypto/rand"
	"math/big"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*_+-"

func RandomString(n int) (string, error) {
	l := int64(len(letters))

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		r, err := rand.Int(rand.Reader, big.NewInt(l))
		if err != nil {
			return "", err
		}

		b[i] = letters[r.Int64()]
	}

	return string(b), nil
}
