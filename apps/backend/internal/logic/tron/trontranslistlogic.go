// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type TronTransListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTronTransListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TronTransListLogic {
	return &TronTransListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TronTransListLogic) TronTransList(req *types.TronTransListReq) (resp *types.TronTransListResp, err error) {
	reqc := &converter.TronTransListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.TronService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find tron trans list failed, err:%v", err)
		return
	}

	resp = &types.TronTransListResp{}
	copier.Copy(resp, res)

	return
}
