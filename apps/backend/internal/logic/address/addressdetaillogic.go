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

type AddressDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressDetailLogic {
	return &AddressDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressDetailLogic) AddressDetail(req *types.AddressDetailReq) (resp *types.AddressItem, err error) {
	// todo: add your logic here and delete this line
	conv, err := l.svcCtx.AddressService.Get(l.ctx, req.Id)
	if err != nil {
		return
	}

	resp = &types.AddressItem{}
	copier.Copy(resp, conv)

	return
}
