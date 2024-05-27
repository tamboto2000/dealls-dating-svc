package entities

type AccountAuthToken struct {
	Entity
	TokenID   string
	AccountID int64
}
