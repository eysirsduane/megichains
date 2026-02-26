// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package fund

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressFundCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundCollectLogic {
	return &AddressFundCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundCollectLogic) AddressFundCollect(req *types.AddressFundCollectReq) (resp *types.AddressFundCollectResp, err error) {
	un := l.ctx.Value("username").(string)
	
	reqc := &converter.AddressFundCollectReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.ChainService.Collect(l.ctx, un, reqc)
	if err != nil {
		logx.Errorf("address fund collect failed, err:%v", err)
		return
	}

	resp = &types.AddressFundCollectResp{}
	copier.Copy(resp, res)

	return
}
