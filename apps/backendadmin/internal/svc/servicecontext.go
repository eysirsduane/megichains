// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megichains/pkg/global"
	"megichains/pkg/service"
)

type ServiceContext struct {
	Config       global.BackendesConfig
	ExcfgService *service.RangeConfigService
	AuthService  *service.AuthService
	UserService  *service.UserService
	OrderService *service.MerchOrderService
	TronService  *service.TronService
}

func NewServiceContext(c global.BackendesConfig, excfg *service.RangeConfigService, auth *service.AuthService, user *service.UserService, order *service.MerchOrderService, tron *service.TronService) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ExcfgService: excfg,
		AuthService:  auth,
		UserService:  user,
		OrderService: order,
		TronService:  tron,
	}
}
