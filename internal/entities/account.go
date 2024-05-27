package entities

type Account struct {
	Entity
	Name     string
	Email    string
	PhoneNum *string
	Password []byte
}
