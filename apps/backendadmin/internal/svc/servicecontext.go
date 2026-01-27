// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megichains/pkg/global"
	"megichains/pkg/service"
)

type ServiceContext struct {
	Config         global.BackendesConfig
	ExcfgService   *service.RangeConfigService
	AuthService    *service.AuthService
	UserService    *service.UserService
	AddressService *service.AddressService
	OrderService   *service.MerchOrderService
	ChainService	*service.ChainService
	TronService    *service.TronService
	EvmService     *service.EvmService
	FundService    *service.FundService
}

func NewServiceContext(c global.BackendesConfig, excfg *service.RangeConfigService, auth *service.AuthService, user *service.UserService, addr *service.AddressService, order *service.MerchOrderService, chain *service.ChainService, tron *service.TronService, evm *service.EvmService, fund *service.FundService) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		ExcfgService:   excfg,
		AuthService:    auth,
		UserService:    user,
		AddressService: addr,
		OrderService:   order,
		ChainService:   chain,
		TronService:    tron,
		EvmService:     evm,
		FundService:    fund,
	}
}
