package entities

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/objects"
)

type Profile struct {
	Entity
	AccountID              int64
	FullName               string
	BirthDate              string
	Gender                 objects.Gender
	RelationNeed           string
	LastEducation          *string
	LastEducationInstitute *string
	ProfilePict            *string
	Headline               *string
	Description            *string
}
