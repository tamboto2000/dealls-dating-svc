package handlers

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/common/accountid"
	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/messages"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/profile"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

func CreateProfile(svc profile.ProfileService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqBody messages.CreateProfileRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			return messages.ResponseError(ctx, errors.New(messages.MsgBadRequest, errors.CodeValidation))
		}

		accId := accountid.GetIDFromCtx(ctx.RequestContext())
		prof, err := profile.NewProfile(
			accId,
			reqBody.BirthDate,
			reqBody.Gender,
			reqBody.RelationNeed,
		)

		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		// headline
		if err := prof.SetHeadline(reqBody.Headline); err != nil {
			return messages.ResponseError(ctx, err)
		}

		// description
		if err := prof.SetDescription(reqBody.Description); err != nil {
			return messages.ResponseError(ctx, err)
		}

		// hobbies
		for _, h := range reqBody.Hobbies {
			if err := prof.AddHobby(h.HobbyName, h.Description); err != nil {
				return messages.ResponseError(ctx, err)
			}
		}

		// interests
		for _, in := range reqBody.Interests {
			if err := prof.AddInterest(in.InterestedIn); err != nil {
				return messages.ResponseError(ctx, err)
			}
		}

		if err := svc.Create(ctx.RequestContext(), prof); err != nil {
			return messages.ResponseError(ctx, err)
		}

		return messages.ResponseSuccess(ctx, nil)
	}
}
