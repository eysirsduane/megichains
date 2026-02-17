// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDetailLogic {
	return &MerchantDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantDetailLogic) MerchantDetail(req *types.MerchantDetailReq) (resp *types.MerchantItem, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.MerchService.GetById(req.Id)
	if err != nil {
		logx.Errorf("find merchant list failed, err:%v", err)
		return
	}

	resp = &types.MerchantItem{}
	copier.Copy(resp, res)

	return
}
