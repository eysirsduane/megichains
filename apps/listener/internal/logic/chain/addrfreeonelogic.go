// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"context"

	"megichains/apps/listener/internal/svc"
	"megichains/apps/listener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddrFreeoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddrFreeoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddrFreeoneLogic {
	return &AddrFreeoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddrFreeoneLogic) AddrFreeone(req *types.AddrFreeoneReq) (resp *types.AddrFreeoneResp, err error) {
	// todo: add your logic here and delete this line

	return
}
