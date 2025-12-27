package service

import (
	"megichains/pkg/entity"

	"gorm.io/gorm"
)

type TronService struct {
	db *gorm.DB
}

func NewTronService(db *gorm.DB) *TronService {
	return &TronService{db: db}
}

func (s *TronService) TransSave(order *entity.TronTransaction) (err error) {
	return s.db.Save(order).Error
}
