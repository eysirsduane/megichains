// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megichains/pkg/global"
	"megichains/pkg/service"
)

type ServiceContext struct {
	Config        global.Config
	ExcfgService  *service.RangeConfigService
	AuthService   *service.AuthService
	UserService   *service.UserService
	ListenService *service.ChainListenService
}

func NewServiceContext(c global.Config, excfg *service.RangeConfigService, auth *service.AuthService, user *service.UserService, listen *service.ChainListenService) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ExcfgService:  excfg,
		AuthService:   auth,
		UserService:   user,
		ListenService: listen,
	}
}
