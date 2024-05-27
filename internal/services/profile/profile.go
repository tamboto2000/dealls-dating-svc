package profile

import (
	"fmt"
	"net/url"
	"regexp"
	"time"
	"unicode"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/internal/objects"
	"github.com/tamboto2000/dealls-dating-svc/pkg/snowid"
)

type Profile struct {
	profile   entities.Profile
	hobbies   []entities.ProfileHobby
	interests []entities.ProfileInterest
}

func NewProfile(accId int64, birth string, gender string, relationNeed string) (*Profile, error) {
	fields := make(errors.Fields)
	// validations.ValidateName(fields, name)
	validateBirth(fields, birth)
	validateRelation(fields, relationNeed)

	if fields.NotEmpty() {
		return nil, errors.NewErrValidation("Profile creation failed: invalid input", fields)
	}

	now := time.Now()
	mainProf := entities.Profile{
		Entity: entities.Entity{
			ID:        snowid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AccountID:    accId,
		BirthDate:    birth,
		Gender:       objects.NewGender(gender),
		RelationNeed: relationNeed,
	}

	profile := Profile{
		profile: mainProf,
	}

	return &profile, nil
}

func (p *Profile) AddHobby(name, desc string) error {
	fields := make(errors.Fields)
	validateHobbyName(fields, name)
	validateHobbyDesc(fields, desc)

	if fields.NotEmpty() {
		return errors.NewErrValidation("Profile creation failed: invalid input", fields)
	}

	now := time.Now()
	p.hobbies = append(p.hobbies, entities.ProfileHobby{
		Entity: entities.Entity{
			ID:        snowid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		HobbyName: name,
		Description: func() *string {
			if desc != "" {
				return &desc
			}

			return nil
		}(),
	})

	return nil
}

func (p *Profile) AddInterest(in string) error {
	fields := make(errors.Fields)
	validateInterest(fields, in)

	now := time.Now()
	p.interests = append(p.interests, entities.ProfileInterest{
		Entity: entities.Entity{
			ID:        snowid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		InterestedIn: in,
	})

	return nil
}

func (p *Profile) SetProfilePict(uri string) error {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid profile picture: %s", err.Error()), errors.CodeValidation)
	}

	return nil
}

func (p *Profile) SetHeadline(headline string) error {
	if headline == "" {
		return nil
	}

	fields := make(errors.Fields)
	validateHeadline(fields, headline)

	if fields.NotEmpty() {
		return errors.NewErrValidation("Invalid profile headline", fields)
	}

	p.profile.Headline = &headline

	return nil
}

func (p *Profile) SetDescription(desc string) error {
	if desc == "" {
		return nil
	}

	field := "description"
	fields := make(errors.Fields)

	if len(desc) > 500 {
		fields.Add(field, "maximum length is 500 characters")
		return errors.NewErrValidation("Invalid description", fields)
	}

	p.profile.Description = &desc

	return nil
}

func (p *Profile) Profile() entities.Profile {
	return p.profile
}

func (p *Profile) Hobbies() []entities.ProfileHobby {
	return p.hobbies
}

func (p *Profile) Interests() []entities.ProfileInterest {
	return p.interests
}

func validateBirth(fields errors.Fields, birth string) {
	field := "birth_date"
	date, err := time.Parse("2006-01-02", birth)
	if err != nil {
		fields.Add(field, "invalid date format")
	}

	age := findAge(date, time.Now())

	if age < 21 {
		fields.Add(field, "minimum age is 21 years old")
	}
}

func validateRelation(fields errors.Fields, relationNeed string) {
	field := "relation_need"

	if len(relationNeed) > 50 {
		fields.Add(field, "maximum length is 50 characters")
	}

	if len(relationNeed) < 5 {
		fields.Add(field, "minimum length is 5 characters")
	}

	for _, r := range relationNeed {
		if unicode.IsControl(r) {
			fields.Add(field, "can only contains letters, numbers, and symbols")
			break
		}
	}
}

func validateHobbyName(fields errors.Fields, name string) {
	field := "hobby_name"

	if len(name) > 50 {
		fields.Add(field, "maximum length is 50 characters")
	}

	if len(name) < 5 {
		fields.Add(field, "minimum length is 5 characters")
	}

	rgx := regexp.MustCompile(`^[a-zA-Z0-9 -]{5,50}$`)
	if !rgx.MatchString(name) {
		fields.Add(field, "invalid character found, allowed characters are a-z, A-Z, 0-9, space and dash")
	}
}

func validateHobbyDesc(fields errors.Fields, desc string) {
	field := "hobby_desc"

	if desc == "" {
		return
	}

	if len(desc) > 500 {
		fields.Add(field, "maximum length is 500 characters")
	}

	if len(desc) < 5 {
		fields.Add(field, "minimum length is 5 characters")
	}
}

func validateInterest(fields errors.Fields, in string) {
	field := "interest"

	if len(in) > 50 {
		fields.Add(field, "maximum length is 50 characters")
	}

	if len(in) < 5 {
		fields.Add(field, "minimum length is 5 characters")
	}
}

func validateHeadline(fields errors.Fields, headline string) {
	field := "headline"
	if headline == "" {
		return
	}

	if len(headline) > 100 {
		fields.Add(field, "maximum length is 100 characters")
	}
}

func findAge(birthdate, today time.Time) int {
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
