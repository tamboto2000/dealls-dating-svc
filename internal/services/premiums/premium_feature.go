package premiums

import (
	"fmt"
	"regexp"
	"time"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/snowid"
)

type PremiumFeature struct {
	feature entities.PremiumFeature
	subs    entities.Subscription
}

func NewPremiumFeature(name, desc string, ts []string) (*PremiumFeature, error) {
	fields := make(errors.Fields)
	validateFeatName(fields, name)
	validateFeatDesc(fields, desc)
	validateFeatTypes(fields, ts)

	if fields.NotEmpty() {
		return nil, errors.NewErrValidation("Failed to create new premium feature: invalid input", fields)
	}

	now := time.Now()
	feature := entities.PremiumFeature{
		Entity: entities.Entity{
			ID:        snowid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		Description: desc,
		Types:       ts,
	}

	premFeat := PremiumFeature{
		feature: feature,
	}

	return &premFeat, nil
}

func (pf *PremiumFeature) Feature() entities.PremiumFeature {
	return pf.feature
}

func (pf *PremiumFeature) SetID(id int64) {
	pf.feature.ID = id
}

func (pf *PremiumFeature) SetTimestamps(created, updated time.Time) {
	pf.feature.CreatedAt = created
	pf.feature.UpdatedAt = updated
}

func (pf *PremiumFeature) Subscribe(accId int64, ty string) error {
	fields := make(errors.Fields)
	validateFeatTypes(fields, []string{ty})
	if fields.NotEmpty() {
		return errors.NewErrValidation("Failed to subscribe to feature: invalid subscription type", fields)
	}

	now := time.Now()
	subs := entities.Subscription{
		Entity: entities.Entity{
			ID:        snowid.Generate(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AccountID:        accId,
		PremiumFeatureID: pf.feature.ID,
		SubsType:         ty,
		ExpiredAt:        expirationTime(now, ty),
	}

	pf.subs = subs

	return nil
}

func (pf PremiumFeature) Subscription() (entities.Subscription, error) {
	if pf.subs.ID == 0 {
		return pf.subs, errors.New("not yet subscribed to this feature", errors.CodeValidation)
	}

	return pf.subs, nil
}

func validateFeatName(fields errors.Fields, name string) {
	field := "name"

	if len(name) > 50 {
		fields.Add(field, "maximum length is 50 characters")
	}

	if len(name) < 3 {
		fields.Add(field, "minimum length is 3 characters")
	}

	rgx := regexp.MustCompile(`^[a-zA-Z0-9 -]{3,50}$`)
	if !rgx.MatchString(name) {
		fields.Add(field, "invalid character found, allowed characters are a-z, A-Z, 0-9, space and dash")
	}
}

func validateFeatTypes(fields errors.Fields, ts []string) {
	field := "types"

	for _, t := range ts {
		switch t {
		case "D", "W", "M", "Y":
			continue
		default:
			fields.Add(field, fmt.Sprintf("unknown type %s, valid types is D for daily, W for weekly, M for monthly, and Y for yearly", t))
		}
	}
}

func validateFeatDesc(fields errors.Fields, desc string) {
	field := "desc"

	if len(desc) > 500 {
		fields.Add(field, "maximum length is 500 characters")
	}

	if len(desc) < 10 {
		fields.Add(field, "minimum length is 10 characters")
	}
}

func expirationTime(start time.Time, ty string) time.Time {
	switch ty {
	case "D":
		return start.AddDate(0, 0, 1)

	case "W":
		return start.AddDate(0, 0, 7)

	case "M":
		return start.AddDate(0, 1, 0)

	case "Y":
		return start.AddDate(1, 0, 0)
	}

	return time.Time{}
}
