// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package solana

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type SolanaTransListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSolanaTransListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SolanaTransListLogic {
	return &SolanaTransListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SolanaTransListLogic) SolanaTransList(req *types.SolanaTransListReq) (resp *types.SolanaTransListResp, err error) {
	reqc := &converter.SolanaTransListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.SolanaService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find solana trans list failed, err:%v", err)
		return
	}

	resp = &types.SolanaTransListResp{}
	copier.Copy(resp, res)

	return
}
