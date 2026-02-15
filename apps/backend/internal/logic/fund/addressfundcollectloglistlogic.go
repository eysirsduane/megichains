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

type AddressFundCollectLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundCollectLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundCollectLogListLogic {
	return &AddressFundCollectLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundCollectLogListLogic) AddressFundCollectLogList(req *types.AddressFundCollectLogListReq) (resp *types.AddressFundCollectLogListResp, err error) {
	reqc := &converter.AddressFundCollectLogListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.FundService.FindCollectLogList(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	resp = &types.AddressFundCollectLogListResp{}
	copier.Copy(resp, res)

	return
}
