// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	atoken, rtoken, err := l.svcCtx.AuthService.RefreshToken(req.RefreshToken, l.svcCtx.Config.Auth.RefreshSecret)
	if err != nil {
		logx.Errorf("auth refresh token failed, err:%v", err)
		return
	}

	resp = &types.RefreshTokenResp{
		Token:        atoken,
		RefreshToken: rtoken,
	}

	return
}
