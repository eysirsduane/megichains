// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package bsc

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BscLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBscLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BscLogListLogic {
	return &BscLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BscLogListLogic) BscLogList(req *types.BscLogListReq) (resp *types.BscLogListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
