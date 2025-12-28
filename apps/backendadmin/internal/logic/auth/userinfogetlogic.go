// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"megichains/apps/backendadmin/internal/svc"
	"megichains/apps/backendadmin/internal/types"

	"github.com/jinzhu/copier"
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
	un := l.ctx.Value("username").(string)
	info, err := l.svcCtx.UserService.Get(un)
	if err != nil {
		logx.Errorf("auth get user info failed, un:%v err:%v", un, err)
		return
	}

	resp = &types.UserInfoGetResp{}
	copier.Copy(resp, info)

	return
}
