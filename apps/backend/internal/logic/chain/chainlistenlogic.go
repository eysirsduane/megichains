// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chain

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/global"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainListenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChainListenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChainListenLogic {
	return &ChainListenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChainListenLogic) ChainListen(req *types.ChainListenReq) (resp *types.ChainListenResp, err error) {
	exist := false
	rkey := global.GetOrderAddressKey(string(req.Chain), req.Receiver, req.Currency)
	l.svcCtx.ListenService.Receivers.Range(func(key, val any) bool {
		if key == rkey {
			logx.Infof("已存在监听地址, chain:%v, receiver:%v, currency:%v", req.Chain, req.Receiver, req.Currency)
			exist = true
			return false
		}

		return true
	})
	if exist {
		err = biz.AlreadyListenThisAddress
		return
	}

	conv := &converter.ChainListenReq{}
	copier.Copy(conv, req)

	go l.svcCtx.ListenService.Listen(conv)

	return
}
