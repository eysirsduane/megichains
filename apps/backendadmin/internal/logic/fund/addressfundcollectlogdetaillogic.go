// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package fund

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressFundCollectLogDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundCollectLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundCollectLogDetailLogic {
	return &AddressFundCollectLogDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundCollectLogDetailLogic) AddressFundCollectLogDetail(req *types.AddressFundCollectLogDetailReq) (resp *types.AddressFundCollectLogItem, err error) {
	item, err := l.svcCtx.FundService.FindCollectLogDetail(l.ctx, req.Id)
	if err != nil {
		return
	}

	resp = &types.AddressFundCollectLogItem{}
	copier.Copy(resp, item)

	return
}
