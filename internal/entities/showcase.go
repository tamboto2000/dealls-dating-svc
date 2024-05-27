package entities

import "github.com/tamboto2000/dealls-dating-svc/internal/objects"

type ShowcaseHistory struct {
	Entity
	AccountID      int64
	ShownProfileID int64
}

type ShowcaseProfile struct {
	Entity
	FullName     string
	BirthDate    string
	Gender       objects.Gender
	RelationNeed string
	Headline     *string
}
