package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type TronService struct {
	db *gorm.DB
}

func NewTronService(db *gorm.DB) *TronService {
	return &TronService{db: db}
}

func (s *TronService) OrderSave(order *entity.TronOrder) (err error) {
	err = s.db.Save(order).Error
	if err != nil {
		logx.Errorf("db tron service save transaction failed, err:%v", err)
		err = biz.EvmOrderSaveFailed
	}

	return
}

func (s *TronService) TransSave(order *entity.TronTransaction) (err error) {
	err = s.db.Save(order).Error
	if err != nil {
		logx.Errorf("db tron service save transaction failed, err:%v", err)
		err = biz.EvmOrderSaveFailed
	}

	return
}
