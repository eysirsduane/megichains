package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SolanaService struct {
	db *gorm.DB
}

func NewSolanaService(db *gorm.DB) *SolanaService {
	return &SolanaService{db: db}
}

func (s *SolanaService) Save(order *entity.SolanaOrder) (err error) {
	err = s.db.Create(order).Error
	if err != nil {
		logx.Errorf("db solana service save transaction failed, err:%v", err)
		err = biz.EvmOrderSaveFailed
	}

	return
}
func (s *SolanaService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db solana service create address failed, err:%v", err)
		err = biz.AddressCreateFailed
	}
	return
}
