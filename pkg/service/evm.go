package service

import (
	"megichains/pkg/entity"

	"gorm.io/gorm"
)

type EvmService struct {
	db *gorm.DB
}

func NewEvmService(db *gorm.DB) *EvmService {
	return &EvmService{db: db}
}

func (s *EvmService) LogSave(blog *entity.EvmLog) (err error) {
	return s.db.Save(blog).Error
}

func (s *EvmService) CreateAddresses(addrs []*entity.Address) (err error) {
	return s.db.Create(addrs).Error
}
