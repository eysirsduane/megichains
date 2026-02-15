// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package fund

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressFundListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundListLogic {
	return &AddressFundListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundListLogic) AddressFundList(req *types.AddressFundListReq) (resp *types.AddressFundListResp, err error) {
	reqc := &converter.AddressFundListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.FundService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	resp = &types.AddressFundListResp{}
	copier.Copy(resp, res)

	return
}
