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

type AddressGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressGenerateLogic {
	return &AddressGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressGenerateLogic) AddressGenerate(req *types.AddressGenerateReq) (resp *types.Response, err error) {
	conv := &converter.ChainAddressCreateReq{}
	copier.Copy(conv, req)

	err = l.svcCtx.AddressService.CreateAddress(conv)

	return
}
