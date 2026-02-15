// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailLogic {
	return &OrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailLogic) OrderDetail(req *types.OrderGetReq) (resp *types.OrderItem, err error) {
	order, err := l.svcCtx.OrderService.Get(req.Id)
	if err != nil {
		logx.Errorf("order get failed, oid:%v, err:%v", req.Id, err)
		return
	}

	resp = &types.OrderItem{}
	copier.Copy(resp, order)

	return
}
