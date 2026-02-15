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

type AddressSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressSaveLogic {
	return &AddressSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressSaveLogic) AddressSave(req *types.AddressItem) (resp *types.Response, err error) {
	reqc := &converter.AddressItem{}
	copier.Copy(reqc, req)

	err = l.svcCtx.AddressService.Save(l.ctx, reqc)
	if err != nil {
		return
	}

	return
}
