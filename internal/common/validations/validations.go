package validations

import (
	"unicode"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
)

func ValidateName(fields errors.Fields, name string) {
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
