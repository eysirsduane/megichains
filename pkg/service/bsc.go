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

func (s *BscService) SaveTransaction(trans *entity.BscTransaction) (err error) {
	err = s.db.Create(trans).Error
	if err != nil {
		logx.Errorf("db bsc service save transaction failed, err:%v", err)
		err = biz.CodeParamsEmpty
	}

	return
}
