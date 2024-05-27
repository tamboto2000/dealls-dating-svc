package objects

import (
	"errors"
	"fmt"
	"strings"
)

var (
	Women = Gender{
		short: "W",
		full:  "Women",
	}

	Men = Gender{
		short: "M",
		full:  "Men",
	}

	Other = Gender{
		short: "O",
		full:  "Other",
	}
)

type Gender struct {
	short string
	full  string
}

func NewGender(what string) Gender {
	return selectGender(what)
}

func (g Gender) Short() string {
	return g.short
}

func (g Gender) Full() string {
	return g.full
}

func (g *Gender) UnmarshalJSON(b []byte) error {
	str := string(b)
	*g = selectGender(str)

	return nil
}

func (g *Gender) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", g.short)), nil
}

func (g *Gender) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("scanning to Gender from non-string")
	}

	*g = selectGender(str)

	return nil
}

func selectGender(str string) Gender {
	var g Gender

	switch strings.ToUpper(str) {
	case "W", "WOMEN":
		g = Women

	case "M", "MEN":
		g = Men

	default:
		g = Other
	}

	return g
}
