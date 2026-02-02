package ping

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	res "go-zero-template/internal/response"

	"go-zero-template/internal/logic/ping"
	"go-zero-template/internal/svc"
	"go-zero-template/internal/types"
)

func PingUserServiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingUserServiceRequest
		if err := httpx.Parse(r, &req); err != nil {
			res.Response(w, nil, res.ParseError)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				res.Response(w, nil, res.InternalServerError)
			}
		}()

		l := ping.NewPingUserServiceLogic(r.Context(), svcCtx)
		resp, err := l.PingUserService(&req)
		res.Response(w, resp, err)
	}
}
