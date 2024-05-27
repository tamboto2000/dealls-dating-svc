package handlers

import (
	"strconv"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/messages"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/showcase"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

func ShowcaseProfiles(svc showcase.ShowcaseService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		profs, err := svc.ShowProfiles(ctx.RequestContext())
		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		var resp []messages.Profile
		for _, p := range profs {
			resp = append(resp, messages.Profile{
				ID:           p.ID,
				Name:         p.FullName,
				BirthDate:    p.BirthDate,
				Gender:       p.Gender.Short(),
				RelationNeed: p.RelationNeed,
				Headline: func(str *string) string {
					if str != nil {
						return *str
					}

					return ""
				}(p.Headline),
			})
		}

		return messages.ResponseSuccess(ctx, resp)
	}
}

func ActionOnShowcase(svc showcase.ShowcaseService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqBody messages.ActionOnShowcaseRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			return messages.ResponseError(ctx, errors.New(messages.MsgBadRequest, errors.CodeValidation))
		}

		fields := make(errors.Fields)
		profIdStr := ctx.PathVal("profile_id")
		if profIdStr == "" {
			fields.Add("profile_id", "can not be empty")

			return messages.ResponseError(ctx, errors.NewErrValidation(messages.MsgBadRequest, fields))
		}

		profId, err := strconv.ParseInt(profIdStr, 10, 64)
		if err != nil {
			fields.Add("profile_id", "must be of type integer")

			return messages.ResponseError(ctx, errors.NewErrValidation(messages.MsgBadRequest, fields))
		}

		if err = svc.Action(ctx.RequestContext(), profId, reqBody.Action); err != nil {
			return messages.ResponseError(ctx, err)
		}

		return messages.ResponseSuccess(ctx, nil)
	}
}
