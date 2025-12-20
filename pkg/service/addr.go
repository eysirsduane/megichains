package service

import (
	"megichains/pkg/entity"

	"gorm.io/gorm"
)

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{db: db}
}

func (s *AddressService) GetAddress(id int64) (addr *entity.Address, err error) {
	err = s.db.Where("id = ?", id).First(&addr).Error
	if err != nil {
		return nil, err
	}

	return
}
