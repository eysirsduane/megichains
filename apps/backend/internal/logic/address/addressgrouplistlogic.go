// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package address

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddressGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressGroupListLogic {
	return &AddressGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressGroupListLogic) AddressGroupList(req *types.AddressGroupListReq) (resp *types.AddressGroupListResp, err error) {
	reqc := &converter.AddressGroupListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.AddressService.GroupFind(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find address list failed, err:%v", err)
		return
	}

	resp = &types.AddressGroupListResp{}
	copier.Copy(resp, res)

	return
}
