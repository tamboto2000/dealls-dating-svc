package messages

import (
	"net/http"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

const (
	StatusError   = "error"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

const (
	MsgBadRequest  = "bad request"
	MsgInternalErr = "internal server error"
)

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  *Error `json:"error,omitempty"`
}

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Fields any    `json:"fields,omitempty"`
}

func ResponseSuccess(ctx *mux.Context, data any) error {
	return ctx.WriteJSON(http.StatusOK, Response{
		Status: StatusSuccess,
		Data:   data,
	})
}

func ResponseError(ctx *mux.Context, err error) error {
	errResp := errToErrResp(err)
	resp := Response{
		Error: &errResp,
	}

	var status int
	switch errResp.Code {
	case errors.CodeValidation:
		status = http.StatusBadRequest
		resp.Status = StatusFailed

	case errors.CodeAlreadyExists:
		status = http.StatusConflict
		resp.Status = StatusFailed

	case errors.CodeInvalidAuth:
		status = http.StatusUnauthorized
		resp.Status = StatusFailed

	case errors.CodeLimitExceeded:
		status = http.StatusTooManyRequests
		resp.Status = StatusFailed

	default:
		status = http.StatusInternalServerError
		resp.Status = StatusError
		resp.Error.Msg = MsgInternalErr
	}

	return ctx.WriteJSON(status, resp)
}

func errToErrResp(err error) Error {
	var errVld errors.ErrValidation
	if errors.As(err, &errVld) {
		return Error{
			Msg:    errVld.Error(),
			Code:   errVld.Code(),
			Fields: errVld.Fields(),
		}
	}

	var errE errors.Err
	if errors.As(err, &errE) {
		return Error{
			Msg:  errE.Error(),
			Code: errE.Code(),
		}
	}

	return Error{
		Msg:  err.Error(),
		Code: errors.CodeInternal,
	}
}
