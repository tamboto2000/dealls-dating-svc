package messages

import "time"

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token     string    `json:"token"`
	TokenID   string    `json:"token_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
