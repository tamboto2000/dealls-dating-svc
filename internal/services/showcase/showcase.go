package showcase

import (
	"context"
	"time"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/accountid"
	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/premiums"
)

type ShowcaseService interface {
	ShowProfiles(ctx context.Context) ([]entities.ShowcaseProfile, error)
	Action(ctx context.Context, profId int64, action string) error
}

type showcaseSvc struct {
	cfg          config.Config
	showcaseRepo ShowcaseRepository
	subsSvc      premiums.SubscriptionService
}

func NewShowcaseService(cfg config.Config, showcaseRepo ShowcaseRepository, subsSvc premiums.SubscriptionService) ShowcaseService {
	return &showcaseSvc{
		cfg:          cfg,
		showcaseRepo: showcaseRepo,
		subsSvc:      subsSvc,
	}
}

func (s *showcaseSvc) ShowProfiles(ctx context.Context) ([]entities.ShowcaseProfile, error) {
	accId := accountid.GetIDFromCtx(ctx)

	if err := s.validateUnlimitedSwipe(ctx); err != nil {
		return nil, err
	}

	return s.showcaseRepo.GetRandomShowcase(ctx, accId, 10)
}

func (s *showcaseSvc) Action(ctx context.Context, profId int64, action string) error {
	accId := accountid.GetIDFromCtx(ctx)
	if err := s.validateUnlimitedSwipe(ctx); err != nil {
		return err
	}

	fields := make(errors.Fields)
	validateAction(fields, action)

	if fields.NotEmpty() {
		return errors.NewErrValidation("Invalid action", fields)
	}

	err := s.showcaseRepo.CreateHistory(ctx, accId, profId, action)

	return err
}

func (s *showcaseSvc) validateUnlimitedSwipe(ctx context.Context) error {
	accId := accountid.GetIDFromCtx(ctx)
	hcount, err := s.showcaseRepo.HistoryCountByDate(ctx, accId, time.Now())
	if err != nil {
		return err
	}

	if hcount >= 10 {
		ulimitFeat := s.cfg.PremiumFeatures.UnlimitedSwipe
		err := s.subsSvc.ValidateSubscription(ctx, ulimitFeat, accId)
		if err != nil {
			if errors.As(err, &errors.ErrNotExists{}) {
				return errors.New(
					"Free daily quota is exceeded, purchase Unlimited Swipe to get unlimited swipes",
					errors.CodeLimitExceeded,
				)
			}

			if errors.As(err, &errors.ErrValidation{}) {
				return errors.New(
					"Renew your Unlimited Swipe subscription to continue swiping without limit",
					errors.CodeLimitExceeded,
				)
			}

			return err
		}

		return nil
	}

	return nil
}

func validateAction(fields errors.Fields, action string) {
	field := "action"
	if action != "L" && action != "P" {
		fields.Add(field, "action must be either L for Like, or P for Pass")
		return
	}
}
