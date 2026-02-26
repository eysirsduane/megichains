// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"megichains/apps/backend/internal/svc"
	"megichains/apps/backend/internal/types"
	"megichains/pkg/converter"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserListGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListGetLogic {
	return &UserListGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListGetLogic) UserListGet(req *types.UserListGetReq) (resp *types.UserListGetResp, err error) {
	reqc := &converter.UserListGetReq{}
	copier.Copy(reqc, req)

	res, err := l.svcCtx.UserService.Find(reqc)
	if err != nil {
		logx.Errorf("find user list failed, err:%v", err)
		return
	}

	resp = &types.UserListGetResp{}
	copier.Copy(resp, res)

	return
}
