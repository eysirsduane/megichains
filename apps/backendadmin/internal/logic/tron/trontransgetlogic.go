// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tron

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"megichains/apps/backendadmin/internal/svc"
)

type TronTransGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTronTransGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TronTransGetLogic {
	return &TronTransGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TronTransGetLogic) TronTransGet() error {
	// todo: add your logic here and delete this line

	return nil
}
