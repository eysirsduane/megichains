package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type BscService struct {
	db *gorm.DB
}

func NewBscService(db *gorm.DB) *BscService {
	return &BscService{db: db}
}

func (s *BscService) SaveLog(blog *entity.BscLog) (err error) {
	err = s.db.Create(blog).Error
	if err != nil {
		logx.Errorf("db bsc service save transaction failed, err:%v", err)
		err = biz.CodeParamsEmpty
	}

	return
}
func (s *BscService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db bsc service create address failed, err:%v", err)
		err = biz.CodeParamsEmpty
	}
	return
}
