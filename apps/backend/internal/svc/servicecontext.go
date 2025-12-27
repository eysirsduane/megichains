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
	UserService   *service.UserService
	AuthService   *service.AuthService
	AddrService   *service.AddressService
	ListenService *service.ChainListenService
}

func NewServiceContext(c global.Config, excfg *service.RangeConfigService, user *service.UserService, auth *service.AuthService, addrservice *service.AddressService, listen *service.ChainListenService) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		ExcfgService:  excfg,
		UserService:   user,
		AuthService:   auth,
		AddrService:   addrservice,
		ListenService: listen,
	}
}
