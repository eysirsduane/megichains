package service

import (
	"context"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/global"

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

func (s *EvmService) CreateAddresses(addrs []*entity.Address) (err error) {
	return s.db.Create(addrs).Error
}

func (s *EvmService) Find(ctx context.Context, req *converter.EvmLogListReq) (resp *converter.RespConverter[entity.EvmLog], err error) {
	db := gorm.G[entity.EvmLog](s.db).Order("updated_at desc")
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Currency != "" {
		db = db.Where("currency = ?", req.Currency)
	}
	if req.TxHash != "" {
		db = db.Where("tx_hash = ?", req.TxHash)
	}
	if req.FromHex != "" {
		db = db.Where("from_hex = ?", req.FromHex)
	}
	if req.ToHex != "" {
		db = db.Where("to_hex = ?", req.ToHex)
	}
	if req.Start > 0 {
		db = db.Where("created_at >= ?", req.Start)
	}
	if req.End > 0 {
		db = db.Where("created_at <= ?", req.End)
	}

	items, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("evm log list find failed, err:%v", err)
		err = biz.OrderFindFailed
		return
	}
	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("evm log list count failed, err:%v", err)
		err = biz.OrderCountFailed
		return
	}

	resp = converter.ConvertToResp(items, req.Current, req.Size, total)
	return
}
