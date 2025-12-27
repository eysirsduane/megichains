// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainAddressCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChainAddressCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChainAddressCreateLogic {
	return &ChainAddressCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChainAddressCreateLogic) ChainAddressCreate(req *types.ChainAddressCreateReq) (resp *types.Response, err error) {
	conv := &converter.ChainAddressCreateReq{}
	copier.Copy(conv, req)

	err = l.svcCtx.AddrService.CreateAddress(conv)

	return
}
