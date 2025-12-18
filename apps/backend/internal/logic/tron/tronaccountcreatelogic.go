// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TronAccountCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTronAccountCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TronAccountCreateLogic {
	return &TronAccountCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TronAccountCreateLogic) TronAccountCreate(req *types.TronAccountCreateReq) (resp *types.TronAccountCreateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
