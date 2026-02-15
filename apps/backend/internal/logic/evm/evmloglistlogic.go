// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package evm

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type EvmLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEvmLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EvmLogListLogic {
	return &EvmLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EvmLogListLogic) EvmLogList(req *types.EvmLogListReq) (resp *types.EvmLogListResp, err error) {
	reqc := &converter.EvmLogListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.EvmService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find evm log list failed, err:%v", err)
		return
	}

	resp = &types.EvmLogListResp{}
	copier.Copy(resp, res)

	return
}
