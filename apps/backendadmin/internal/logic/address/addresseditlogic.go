// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package address

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressEditLogic {
	return &AddressEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressEditLogic) AddressEdit(req *types.AddressItem) (err error) {
	reqc := &converter.AddressItem{}
	copier.Copy(reqc, req)

	err = l.svcCtx.AddressService.Edit(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	return
}
