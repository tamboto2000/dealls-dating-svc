package account

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/bcrypt"
	"github.com/tamboto2000/dealls-dating-svc/pkg/snowid"
)

type Account struct {
	acc *entities.Account
}

func NewAccount(name, email, phoneNum, pwd string, pwdCfg config.Password) (*Account, error) {
	// validation steps
	fields := make(errors.Fields)
	validateName(fields, name)
	validateEmail(fields, email)
	validatePwd(pwdCfg, fields, pwd)
	if phoneNum != "" {
		validatePhone(fields, phoneNum)
	}

	if fields.NotEmpty() {
		return nil, errors.NewErrValidation("Account creation failed: invalid data input", fields)
	}

	pwd = fmt.Sprintf("%s%s", pwd, pwdCfg.Salt)
	hashedPwd, err := bcrypt.Hash([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	acc := Account{
		acc: &entities.Account{
			Entity: entities.Entity{
				ID:        snowid.Generate(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:  name,
			Email: email,
			PhoneNum: func(phoneNum string) *string {
				if phoneNum != "" {
					return &phoneNum
				}

				return nil
			}(phoneNum),
			Password: hashedPwd,
		},
	}

	return &acc, nil
}

func (acc Account) Email() string {
	return acc.acc.Email
}

func (acc Account) PhoneNum() string {
	var phoneNum string
	if acc.acc.PhoneNum != nil {
		phoneNum = *acc.acc.PhoneNum
	}

	return phoneNum
}

// validateName validates user name.
// name len can not be more than 100 char and
// minimum 1 char.
// Name can not contains control character
func validateName(fields errors.Fields, name string) {
	field := "name"

	if name == "" {
		fields.Add(field, "name can not be empty")
		return
	}

	// check name length
	if len(name) > 100 {
		fields.Add(field, "name length is more than 100 characters")
	}

	if len(name) == 0 {
		fields.Add(field, "name length is less than 1 character")
	}

	// check if name contains control characters
	for _, r := range name {
		if unicode.IsControl(r) {
			fields.Add(field, "name contains illegal characters")
			break
		}
	}
}

func validateEmail(fields errors.Fields, email string) {
	field := "email"

	if email == "" {
		fields.Add(field, "email can not be empty")
		return
	}

	if strings.ContainsAny(email, " ") {
		fields.Add(field, "invalid email address")
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		fields.Add(field, "invalid email address")
	}
}

func validatePhone(fields errors.Fields, phone string) {
	field := "phone"

	if len(phone) > 13 {
		fields.Add(field, "phone number length is more than 13 characters")
	}

	if len(phone) < 10 {
		fields.Add(field, "phone number length is less than 10 characters")
	}
}

func validatePwd(pwdCfg config.Password, fields errors.Fields, pwd string) {
	field := "password"

	if pwd == "" {
		fields.Add(field, "password can not be empty")
		return
	}

	uppercase := false
	lowercase := false
	number := false
	specialChar := false
	letters := 0

	for _, r := range pwd {
		letters++

		if unicode.IsUpper(r) {
			uppercase = true
			continue
		}

		if unicode.IsLower(r) {
			lowercase = true
			continue
		}

		if unicode.IsNumber(r) {
			number = true
			continue
		}

		if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			specialChar = true
		}
	}

	if !uppercase || !lowercase || !number || !specialChar {
		fields.Add(field, "password must contains at least 1 upper case, 1 lower case, 1 numeric, and 1 special character")
	}

	if letters < pwdCfg.MinLength {
		fields.Add(field, fmt.Sprintf("password length is less than %d characters", pwdCfg.MinLength))
	}

	if letters > pwdCfg.MaxLength {
		fields.Add(field, fmt.Sprintf("password length is more than %d characters", pwdCfg.MaxLength))
	}
}
