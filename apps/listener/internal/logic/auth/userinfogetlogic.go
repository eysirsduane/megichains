// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"megichains/apps/listener/internal/svc"
	"megichains/apps/listener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoGetLogic {
	return &UserInfoGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoGetLogic) UserInfoGet(req *types.UserInfoGetReq) (resp *types.UserInfoGetResp, err error) {
	// todo: add your logic here and delete this line

	return
}
