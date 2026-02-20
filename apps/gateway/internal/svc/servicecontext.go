// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megichains/apps/gateway/internal/middleware"
	"megichains/pkg/global"
	"megichains/pkg/service"
)

type ServiceContext struct {
	Config                        global.BackendesConfig
	ApiAccessPermissionMiddleware *middleware.ListenMiddleware
	ExcfgService                  *service.RangeConfigService
	UserService                   *service.UserService
	AuthService                   *service.AuthService
	AddrService                   *service.AddressService
	ListenService                 *service.ListenService
}

func NewServiceContext(c global.BackendesConfig, apimiddle *middleware.ListenMiddleware, excfg *service.RangeConfigService, user *service.UserService, auth *service.AuthService, addrservice *service.AddressService, listen *service.ListenService) *ServiceContext {
	return &ServiceContext{
		Config:                        c,
		ApiAccessPermissionMiddleware: apimiddle,
		ExcfgService:                  excfg,
		UserService:                   user,
		AuthService:                   auth,
		AddrService:                   addrservice,
		ListenService:                 listen,
	}
}
