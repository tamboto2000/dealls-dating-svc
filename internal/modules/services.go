package modules

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/account"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/premiums"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/profile"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/showcase"
)

type Services struct {
	accSvc      account.AccountService
	profileSvc  profile.ProfileService
	showcaseSvc showcase.ShowcaseService
	subsSvc     premiums.SubscriptionService
}

func NewServices(ctx context.Context, cfg config.Config, in *Infra, comps *Components) *Services {
	// account service
	accRepo := account.NewAccountRepo(in.DB())
	accsvc := account.NewAccountService(cfg, accRepo, comps.Mailer())

	// profile service
	profileRepo := profile.NewProfileRepository(in.DB())
	profileSvc := profile.NewProfileService(profileRepo)

	// subscription service
	premFeatRepo := premiums.NewPremiumFeatureRepository(in.DB())
	subsRepo := premiums.NewSubscriptionRepository(in.DB(), in.Cache())
	subsSvc := premiums.NewSubscriptionService(premFeatRepo, subsRepo)

	// showcase service
	showcaseRepo := showcase.NewShowcaseRepository(in.DB(), in.Cache())
	showcaseSvc := showcase.NewShowcaseService(cfg, showcaseRepo, subsSvc)

	// register services
	svcs := Services{
		accSvc:      accsvc,
		profileSvc:  profileSvc,
		showcaseSvc: showcaseSvc,
		subsSvc:     subsSvc,
	}

	return &svcs
}

func (s *Services) AccountService() account.AccountService {
	return s.accSvc
}

func (s *Services) ProfileService() profile.ProfileService {
	return s.profileSvc
}

func (s *Services) ShowcaseService() showcase.ShowcaseService {
	return s.showcaseSvc
}

func (s *Services) SubscriptionService() premiums.SubscriptionService {
	return s.subsSvc
}
