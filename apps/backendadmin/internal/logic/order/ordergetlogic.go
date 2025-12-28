// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderGetLogic {
	return &OrderGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderGetLogic) OrderGet(req *types.OrderGetReq) (resp *types.OrderItem, err error) {
	// todo: add your logic here and delete this line
	order, err := l.svcCtx.OrderService.Get(req.Id)
	if err != nil {
		logx.Errorf("order get failed, oid:%v, err:%v", req.Id, err)
		return
	}

	resp = &types.OrderItem{}
	copier.Copy(resp, order)

	return
}
