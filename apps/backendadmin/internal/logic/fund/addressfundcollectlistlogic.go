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

type AddressFundCollectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundCollectListLogic {
	return &AddressFundCollectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundCollectListLogic) AddressFundCollectList(req *types.AddressFundCollectListReq) (resp *types.AddressFundCollectListResp, err error) {
	reqc := &converter.AddressFundCollectListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.FundService.FindCollectLogList(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	resp = &types.AddressFundCollectListResp{}
	copier.Copy(resp, res)

	return
}
