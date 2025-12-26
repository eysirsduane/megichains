package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
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

func (s *EvmService) OrderFirstOrCreate(order *entity.EvmOrder) (err error) {
	return s.db.FirstOrCreate(order).Error
}

func (s *EvmService) OrderSave(order *entity.EvmOrder) (err error) {
	return s.db.Save(order).Error
}

func (s *EvmService) GetOrder(moid string) (order *entity.EvmOrder, err error) {
	order = &entity.EvmOrder{}
	err = s.db.Model(&entity.EvmOrder{}).Where("merch_order_id = ?", moid).First(order).Error
	return
}

func (s *EvmService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db evm service create address failed, err:%v", err)
		err = biz.AddressCreateFailed
	}
	return
}
