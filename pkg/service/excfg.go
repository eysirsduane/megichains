package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RangeConfigService struct {
	db *gorm.DB
}

func NewRangeConfigService(db *gorm.DB) *RangeConfigService {
	return &RangeConfigService{db: db}
}

func (s *RangeConfigService) GetConfig(typo string, rfrom, rto int64) (cfg *entity.RangeConfig, err error) {
	cfg = &entity.RangeConfig{}

	err = s.db.Model(&entity.RangeConfig{}).Where("typo = ? and  range_from >= ? and range_to <= ?", typo, rfrom, rto).First(cfg).Error
	if err != nil {
		logx.Errorf("db range config get failed, typo:%v, rfrom:%v, rto:%v, err:%v", typo, rfrom, rto, err)
		err = biz.RangeConfigGetFailed
		return
	}

	if cfg.Value == 0 {
		logx.Errorf("db range config value is 0, typo:%v, rfrom:%v, rto:%v, err:%v", typo, rfrom, rto, err)
		err = biz.RangeConfigValueInvalid
		return
	}
	if typo == string(global.ExchangeTypoUsdt2Trx) {
		if cfg.Value >= 1 {
			logx.Errorf("db range config usdt2trx value invalid, value:%v", cfg.Value)
			err = biz.RangeConfigValueInvalid
			return
		}
	}

	return
}
