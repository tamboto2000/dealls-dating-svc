package jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSigned(t *testing.T) {
	claims := make(Claims)
	claims.SetIssuedAt(time.Now())
	claims.SetExpiresAt(time.Now().Add(24 * time.Hour))
	claims.Set("custom", "Hello, World!")

	tokenStr, err := NewSigned(claims, HS256, []byte("12345678"))
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println("token:", tokenStr)

	var decClaims Claims
	if err := Decode(tokenStr, decClaims, []byte("12345678")); err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("decoded: %#v\n", decClaims)
}
