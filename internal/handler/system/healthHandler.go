package system

import (
	"net/http"

	"go-zero-template/internal/logic/system"
	res "go-zero-template/internal/response"
	"go-zero-template/internal/svc"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				res.Response(w, nil, res.InternalServerError)
			}
		}()

		l := system.NewHealthLogic(r.Context(), svcCtx)
		resp, err := l.Health()
		res.Response(w, resp, err)
	}
}
