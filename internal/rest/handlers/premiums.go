package handlers

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/messages"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/premiums"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

func PremiumFeatureList(svc premiums.SubscriptionService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		feats, err := svc.GetAllFeatures(ctx.RequestContext())
		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		var respData messages.PremiumFeaturesResponse

		for _, f := range feats {
			respData = append(respData, messages.PremiumFeature{
				ID:          f.ID,
				Name:        f.Name,
				Description: f.Description,
				Types:       f.Types,
			})
		}

		return messages.ResponseSuccess(ctx, respData)
	}
}
