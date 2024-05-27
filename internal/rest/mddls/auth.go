package mddls

import (
	"context"
	"regexp"
	"strings"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/messages"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/account"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

var errInvalidAuth = errors.New("Invalid authorization token", errors.CodeInvalidAuth)

func Auth(svc account.AccountService) mux.Middleware {
	return func(hf mux.HandleFunc) mux.HandleFunc {
		return func(ctx *mux.Context) error {
			// exclusion
			rgx := regexp.MustCompile(`(\/register)|(\/signin)`)
			if rgx.MatchString(ctx.Path()) {
				return hf(ctx)
			}

			bearer := ctx.GetHeader("Authorization")
			splits := strings.Split(bearer, " ")
			if len(splits) < 2 {
				return messages.ResponseError(ctx, errInvalidAuth)
			}

			if splits[0] != "Bearer" {
				return messages.ResponseError(ctx, errInvalidAuth)
			}

			accAuth, err := svc.Authorization(ctx.RequestContext(), splits[1])
			if err != nil {
				logger.Error(err.Error())

				return messages.ResponseError(ctx, err)
			}

			accId, err := accAuth.AccountID()
			if err != nil {
				logger.Error(err.Error())

				return messages.ResponseError(ctx, errors.New(messages.MsgInternalErr, errors.CodeInternal))
			}

			// TODO: use custom key type for context value
			reqCtx := context.WithValue(ctx.RequestContext(), "account_id", accId)
			ctx.SetRequestContext(reqCtx)

			return hf(ctx)
		}
	}
}
