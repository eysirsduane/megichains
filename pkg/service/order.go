package service

import (
	"megichains/pkg/entity"

	"gorm.io/gorm"
)

type MerchOrderService struct {
	db *gorm.DB
}

func NewMerchOrderService(db *gorm.DB) *MerchOrderService {
	return &MerchOrderService{db: db}
}

func (s *MerchOrderService) Save(order *entity.MerchOrder) (err error) {
	return s.db.Save(order).Error
}

func (s *MerchOrderService) Get(moid string) (order *entity.MerchOrder, err error) {
	order = &entity.MerchOrder{}
	err = s.db.Model(&entity.MerchOrder{}).Where("merch_order_id = ?", moid).First(order).Error
	return
}
