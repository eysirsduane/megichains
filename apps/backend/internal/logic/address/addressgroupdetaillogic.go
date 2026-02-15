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

type AddressGroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressGroupDetailLogic {
	return &AddressGroupDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressGroupDetailLogic) AddressGroupDetail(req *types.AddressGroupDetailReq) (resp *types.AddressGroupItem, err error) {
	// todo: add your logic here and delete this line
	group, err := l.svcCtx.AddressService.GroupGet(l.ctx, req.Id)
	if err != nil {
		return
	}

	resp = &types.AddressGroupItem{}
	copier.Copy(resp, group)

	return
}
