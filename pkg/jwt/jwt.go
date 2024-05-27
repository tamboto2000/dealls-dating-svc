package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrSignatureInvalid = errors.New("signature is invalid")
)

type SignMethod int

const (
	HS256 SignMethod = iota
)

func NewSigned(claims Claims, signMethod SignMethod, key any) (string, error) {
	var signer jwt.SigningMethod

	switch signMethod {
	case HS256:
		signer = jwt.SigningMethodHS256
	default:
		signer = jwt.SigningMethodNone
	}

	claimsw := claimsWrap(claims)
	token := jwt.NewWithClaims(signer, claimsw)

	return token.SignedString(key)
}

func Decode(tokenStr string, claims Claims, key any) error {
	claimsw := claimsWrap(claims)
	token, err := jwt.ParseWithClaims(tokenStr, &claimsw, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if err.Error() == "token signature is invalid: signature is invalid" {
			return ErrSignatureInvalid
		}

		return err
	}

	if !token.Valid {
		return ErrTokenInvalid
	}

	return nil
}

type claimsWrap Claims

func (cw claimsWrap) GetExpirationTime() (*jwt.NumericDate, error) {
	t, _ := cw["exp"].(*jwt.NumericDate)

	return t, nil
}

func (cw claimsWrap) GetIssuedAt() (*jwt.NumericDate, error) {
	t, _ := cw["iat"].(time.Time)

	return jwt.NewNumericDate(t), nil
}

func (cw claimsWrap) GetNotBefore() (*jwt.NumericDate, error) {
	t, _ := cw["nbf"].(time.Time)

	return jwt.NewNumericDate(t), nil
}

func (cw claimsWrap) GetIssuer() (string, error) {
	s, _ := cw["iss"].(string)

	return s, nil
}
func (cw claimsWrap) GetSubject() (string, error) {
	s, _ := cw["sub"].(string)

	return s, nil
}

func (cw claimsWrap) GetAudience() (jwt.ClaimStrings, error) {
	s, _ := cw["aud"].([]string)

	return s, nil
}
