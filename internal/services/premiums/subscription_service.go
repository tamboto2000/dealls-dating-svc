package premiums

import (
	"context"
	"time"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
)

type SubscriptionService interface {
	GetAllFeatures(ctx context.Context) ([]entities.PremiumFeature, error)
	ValidateSubscription(ctx context.Context, featId, accId int64) error
}

type subsSvc struct {
	premFeatRepo PremiumFeatureRepository
	subsRepo     SubscriptionRepository
}

func NewSubscriptionService(premFeatRepo PremiumFeatureRepository, subsRepo SubscriptionRepository) SubscriptionService {
	return &subsSvc{
		premFeatRepo: premFeatRepo,
		subsRepo:     subsRepo,
	}
}

func (s *subsSvc) GetAllFeatures(ctx context.Context) ([]entities.PremiumFeature, error) {
	return s.premFeatRepo.GetAll(ctx)
}

func (s *subsSvc) ValidateSubscription(ctx context.Context, featId, accId int64) error {
	subs, err := s.subsRepo.GetByAccountID(ctx, featId, accId)
	if err != nil {
		if err == errSubsNotExists {
			return errors.NewErrNotExists("user has not subscribed to this feature")
		}

		return err
	}

	if !subs.ExpiredAt.Before(time.Now()) {
		return errors.NewErrValidation("subscription expired", nil)
	}

	return nil
}

func (s *subsSvc) Subscribe(ctx context.Context, featId, accId int64, ty string) (entities.Subscription, error) {
	feat, err := s.premFeatRepo.Get(ctx, featId)
	if err != nil {
		return entities.Subscription{}, err
	}

	var subs entities.Subscription

	tyFound := false
	for _, t := range feat.Types {
		if ty == t {
			tyFound = true
			break
		}
	}

	if !tyFound {
		switch ty {
		case "D":
			return subs, errors.New("Daily subscription is not supported by this feature", errors.CodeValidation)

		case "W":
			return subs, errors.New("Weekly subscription is not supported by this feature", errors.CodeValidation)

		case "M":
			return subs, errors.New("Monthly subscription is not supported by this feature", errors.CodeValidation)

		case "Y":
			return subs, errors.New("Yearly subscription is not supported by this feature", errors.CodeValidation)
		}

		fields := make(errors.Fields)
		fields.Add("subs_type", "Subscription model is not supported")

		return subs, errors.NewErrValidation("Failed to subscribe to feature: invalid subscription type", fields)
	}

	

	return entities.Subscription{}, nil
}
