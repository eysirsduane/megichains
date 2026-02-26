// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"context"

	"megichains/apps/listener/internal/svc"
	"megichains/apps/listener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChainListensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChainListensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChainListensLogic {
	return &ChainListensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChainListensLogic) ChainListens() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	l.svcCtx.ListenService.ListenMany()

	return
}
