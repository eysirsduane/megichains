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

type AddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressListLogic {
	return &AddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressListLogic) AddressList(req *types.AddressListReq) (resp *types.AddressListResp, err error) {
	reqc := &converter.AddressListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.AddressService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	resp = &types.AddressListResp{}
	copier.Copy(resp, res)

	return
}
