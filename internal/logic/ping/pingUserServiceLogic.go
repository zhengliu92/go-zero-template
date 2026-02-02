// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ping

import (
	"context"
	"errors"

	"go-zero-template/internal/middleware"
	"go-zero-template/internal/svc"
	"go-zero-template/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingUserServiceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingUserServiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingUserServiceLogic {
	return &PingUserServiceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingUserServiceLogic) PingUserService(req *types.PingUserServiceRequest) (resp *types.PingUserServiceResponse, err error) {
	user, ok := middleware.GetUserFromContext(l.ctx)
	if !ok {
		return nil, errors.New("user not found")
	}

	return &types.PingUserServiceResponse{
		ID:            user.ID,
		Name:          user.Name,
		SAPEmployeeID: user.SAPEmployeeID,
	}, nil
}
