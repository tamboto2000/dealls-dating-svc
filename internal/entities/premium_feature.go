package entities

import "time"

type PremiumFeature struct {
	Entity
	Name        string
	Description string
	Types       []string
}

type Subscription struct {
	Entity
	AccountID        int64
	PremiumFeatureID int64
	SubsType         string
	Status           string
	ExpiredAt        time.Time
}
