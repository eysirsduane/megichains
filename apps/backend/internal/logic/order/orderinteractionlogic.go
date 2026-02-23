// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderInteractionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderInteractionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderInteractionLogic {
	return &OrderInteractionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderInteractionLogic) OrderInteraction(req *types.OrderInteractionReq) (resp *types.OrderInteractionResp, err error) {
	reqc := &converter.OrderInteractionReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.OrderService.Interaction(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find order failed, err:%v", err)
		return
	}

	resp = &types.OrderInteractionResp{}
	copier.Copy(resp, res)

	return
}
