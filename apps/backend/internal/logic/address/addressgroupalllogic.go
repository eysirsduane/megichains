// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package address

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressGroupAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressGroupAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressGroupAllLogic {
	return &AddressGroupAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressGroupAllLogic) AddressGroupAll() (resp *types.AddressGroupAllResp, err error) {
	res, err := l.svcCtx.AddressService.GroupAll(l.ctx)
	if err != nil {
		logx.Errorf("find address group list failed, err:%v", err)
		return
	}

	resp = &types.AddressGroupAllResp{}
	copier.Copy(resp, res)

	return
}
