package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type EvmService struct {
	db *gorm.DB
}

func NewEvmService(db *gorm.DB) *EvmService {
	return &EvmService{db: db}
}

func (s *EvmService) SaveLog(blog *entity.EvmOrder) (err error) {
	err = s.db.Create(blog).Error
	if err != nil {
		logx.Errorf("db evm service save transaction failed, err:%v", err)
		err = biz.EvmOrderSaveFailed
	}

	return
}
func (s *EvmService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db evm service create address failed, err:%v", err)
		err = biz.AddressCreateFailed
	}
	return
}
