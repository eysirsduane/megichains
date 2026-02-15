// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"megichains/apps/backend/internal/svc"
)

type TronTransDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTronTransDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TronTransDetailLogic {
	return &TronTransDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TronTransDetailLogic) TronTransDetail() error {
	// todo: add your logic here and delete this line

	return nil
}
