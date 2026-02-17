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

type MerchantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListLogic {
	return &MerchantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListLogic) MerchantList(req *types.MerchantListReq) (resp *types.MerchantListResp, err error) {
	reqc := &converter.MerchantListReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.MerchService.Find(l.ctx, reqc)
	if err != nil {
		logx.Errorf("find merchant list failed, err:%v", err)
		return
	}

	resp = &types.MerchantListResp{}
	copier.Copy(resp, res)

	return
}
