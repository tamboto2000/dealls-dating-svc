package handlers

import (
	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest/messages"
	"github.com/tamboto2000/dealls-dating-svc/internal/services/account"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

func RegisterAccount(cfg config.Config, svc account.AccountService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqAcc messages.RegisterAccountRequest
		if err := ctx.BindJSON(&reqAcc); err != nil {
			return messages.ResponseError(ctx, errors.New(messages.MsgBadRequest, errors.CodeValidation))
		}

		pwdCfg := cfg.Auth.Password
		acc, err := account.NewAccount(reqAcc.Name, reqAcc.Email, reqAcc.Phone, reqAcc.Password, pwdCfg)
		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		if err := svc.CreateAccount(ctx.RequestContext(), acc); err != nil {
			return messages.ResponseError(ctx, err)
		}

		return messages.ResponseSuccess(ctx, nil)
	}
}

func SignIn(svc account.AccountService) mux.HandleFunc {
	return func(ctx *mux.Context) error {
		var reqBody messages.SignInRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			return messages.ResponseError(ctx, errors.New(messages.MsgBadRequest, errors.CodeValidation))
		}

		accAuth, err := svc.SignIn(ctx.RequestContext(), reqBody.Username, reqBody.Password)
		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		token, err := accAuth.Token()
		if err != nil {
			return messages.ResponseError(ctx, err)
		}

		tokenId, _ := accAuth.TokenID()
		expiredAt, _ := accAuth.TokenExpiration()

		resp := messages.SignInResponse{
			Token:     token,
			TokenID:   tokenId,
			ExpiredAt: expiredAt,
		}

		return messages.ResponseSuccess(ctx, resp)
	}
}
