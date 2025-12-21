package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type EthService struct {
	db *gorm.DB
}

func NewEthService(db *gorm.DB) *EthService {
	return &EthService{db: db}
}

func (s *EthService) SaveLog(blog *entity.EthOrder) (err error) {
	err = s.db.Create(blog).Error
	if err != nil {
		logx.Errorf("db eth service save transaction failed, err:%v", err)
		err = biz.CodeParamsEmpty
	}

	return
}
func (s *EthService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db eth service create address failed, err:%v", err)
		err = biz.CodeParamsEmpty
	}
	return
}
