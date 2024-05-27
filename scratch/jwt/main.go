package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	ss, err := generateJwt()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fmt.Println("token:", ss)

	claims, err := decodejwt(ss)
	if err != nil {
		slog.Error(err.Error())
	}

	fmt.Printf("parsed claims: %#v\n", claims)
}

const key = "mysupersecretkey"

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJwt() (string, error) {
	claims := CustomClaims{
		UserID: 1234567890,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func decodejwt(ss string) (CustomClaims, error) {
	var claims CustomClaims
	token, err := jwt.ParseWithClaims(ss, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			
		}

		return claims, err
	}

	if !token.Valid {
		return claims, errors.New("token is invalid")
	}

	return claims, nil
}
