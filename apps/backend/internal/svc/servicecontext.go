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
	MerchService   *service.MerchService
	AddressService *service.AddressService
	OrderService   *service.MerchOrderService
	ChainService   *service.ChainService
	TronService    *service.TronService
	EvmService     *service.EvmService
	FundService    *service.FundService
	SolanaService  *service.SolanaService
}

func NewServiceContext(c global.BackendesConfig, excfg *service.RangeConfigService, auth *service.AuthService, user *service.UserService, merch *service.MerchService, addr *service.AddressService, order *service.MerchOrderService, chain *service.ChainService, tron *service.TronService, evm *service.EvmService, fund *service.FundService, solana *service.SolanaService) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		ExcfgService:   excfg,
		AuthService:    auth,
		UserService:    user,
		MerchService:   merch,
		AddressService: addr,
		OrderService:   order,
		ChainService:   chain,
		TronService:    tron,
		EvmService:     evm,
		FundService:    fund,
		SolanaService:  solana,
	}
}
