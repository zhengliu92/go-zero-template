package response

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Response(w http.ResponseWriter, resp any, err error) {
	var res Result
	if err != nil {
		var e *Error
		switch {
		case errors.As(err, &e):
			res.Code = e.Code
			res.Msg = e.Msg
		default:
			res.Code = http.StatusInternalServerError
			res.Msg = err.Error()
		}
	} else {
		res.Code = http.StatusOK
		res.Msg = SUCCESS
		res.Data = resp
	}
	httpx.OkJson(w, res)
}
