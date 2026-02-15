// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package fund

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressFundStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressFundStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressFundStatisticsLogic {
	return &AddressFundStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressFundStatisticsLogic) AddressFundStatistics(req *types.AddressFundStatisticsReq) (resp *types.AddressFundStatisticsResp, err error) {
	res, err := l.svcCtx.FundService.Statistics(l.ctx)
	if err != nil {
		logx.Errorf("address fund statistics failed, err:%v", err)
		return
	}

	resp = &types.AddressFundStatisticsResp{}
	copier.Copy(resp, res)

	return
}
