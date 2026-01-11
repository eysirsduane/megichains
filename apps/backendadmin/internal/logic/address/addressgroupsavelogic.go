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

type AddressGroupSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressGroupSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressGroupSaveLogic {
	return &AddressGroupSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressGroupSaveLogic) AddressGroupSave(req *types.AddressGroupItem) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	conv := &converter.AddressGroupItem{}
	copier.Copy(conv, req)

	err = l.svcCtx.AddressService.GroupSave(l.ctx, conv)
	if err != nil {
		return
	}

	return
}
