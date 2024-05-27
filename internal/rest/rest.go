package rest

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/modules"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/handlers"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/mddls"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

func RegisterREST(cfg config.Config, r *mux.Router, svcs *modules.Services) {
	rv1 := r.SubRouter("/api/v1")
	rv1.Use(mddls.Auth(svcs.AccountService()))

	// auth
	rv1.Post("/accounts/register", handlers.RegisterAccount(cfg, svcs.AccountService()))
	rv1.Post("/accounts/signin", handlers.SignIn(svcs.AccountService()))

	// profile
	rv1.Post("/profiles", handlers.CreateProfile(svcs.ProfileService()))

	// showcase
	rv1.Get("/showcases", handlers.ShowcaseProfiles(svcs.ShowcaseService()))
	rv1.Post("/showcases/{profile_id}/action", handlers.ActionOnShowcase(svcs.ShowcaseService()))

	// premium features and subscriptions
	rv1.Get("/premium_features", handlers.PremiumFeatureList(svcs.SubscriptionService()))
}
