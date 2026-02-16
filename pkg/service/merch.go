package service

import (
	"megichains/pkg/entity"

	"gorm.io/gorm"
)

type MerchService struct {
	db *gorm.DB
}

func NewMerchService(db *gorm.DB) *MerchService {
	return &MerchService{db: db}
}

func (s *MerchService) Get(account string) (merch *entity.Merchant, err error) {
	merch = &entity.Merchant{}
	err = s.db.Model(&entity.Merchant{}).Where("merchant_account = ?", account).Error

	return
}
