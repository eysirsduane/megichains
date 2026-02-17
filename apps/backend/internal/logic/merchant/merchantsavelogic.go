// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantSaveLogic {
	return &MerchantSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantSaveLogic) MerchantSave(req *types.MerchantItem) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	reqc := &converter.MerchantItem{}
	copier.Copy(reqc, req)

	err = l.svcCtx.MerchService.Save(l.ctx, reqc)
	if err != nil {
		return
	}

	return
}
