// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package bsc

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BscTransactionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBscTransactionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BscTransactionListLogic {
	return &BscTransactionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BscTransactionListLogic) BscTransactionList(req *types.BscTransactionListReq) (resp *types.BscTransactionListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
