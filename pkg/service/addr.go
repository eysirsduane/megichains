package service

import (
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
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
		logx.Errorf("db get address failed, id:%v, err:%v", id, err)
		return
	}

	return
}
