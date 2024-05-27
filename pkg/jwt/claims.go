package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims map[string]any

func (c Claims) Issuer() string {
	s, _ := c["iss"].(string)

	return s
}

func (c Claims) Subject() string {
	s, _ := c["sub"].(string)

	return s
}

func (c Claims) Audience() []string {
	s, _ := c["aud"].(jwt.ClaimStrings)

	return s
}

func (c Claims) ExpiresAt() time.Time {
	v := c["exp"]
	switch t := v.(type) {
	case int64:
		return time.Unix(t, 0)

	case float64:
		return time.Unix(int64(t), 0)
	}

	return time.Time{}
}

func (c Claims) NotBefore() time.Time {
	v := c["nbf"]
	switch t := v.(type) {
	case int64:
		return time.Unix(t, 0)

	case float64:
		return time.Unix(int64(t), 0)
	}

	return time.Time{}
}

func (c Claims) IssuedAt() time.Time {
	v := c["iat"]
	switch t := v.(type) {
	case int64:
		return time.Unix(t, 0)

	case float64:
		return time.Unix(int64(t), 0)
	}

	return time.Time{}
}

func (c Claims) ID() string {
	s, _ := c["jti"].(string)

	return s
}

func (c Claims) Val(key string) any {
	return c[key]
}

func (c Claims) SetIssuer(iss string) {
	c["iss"] = iss
}

func (c Claims) SetSubject(sub string) {
	c["sub"] = sub
}

func (c Claims) SetAudience(aud []string) {
	c["aud"] = jwt.ClaimStrings(aud)
}

func (c Claims) SetExpiresAt(exp time.Time) {
	c["exp"] = jwt.NewNumericDate(exp)
}

func (c Claims) SetNotBefore(nbf time.Time) {
	c["nbf"] = jwt.NewNumericDate(nbf)
}

func (c Claims) SetIssuedAt(iat time.Time) {
	c["iat"] = jwt.NewNumericDate(iat)
}

func (c Claims) SetID(jti string) {
	c["jti"] = jti
}

func (c Claims) Set(key string, val any) {
	c[key] = val
}
