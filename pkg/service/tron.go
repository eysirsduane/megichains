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

type TronService struct {
	db *gorm.DB
}

func NewTronService(db *gorm.DB) *TronService {
	return &TronService{db: db}
}

func (s *TronService) TransSave(order *entity.TronTransaction) (err error) {
	return s.db.Save(order).Error
}

func (s *TronService) Find(ctx context.Context, req *converter.TronTransListReq) (resp *converter.RespConverter[entity.TronTransaction], err error) {
	db := gorm.G[entity.TronTransaction](s.db).Order("updated_at desc")
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if req.Currency != "" {
		db = db.Where("currency = ?", req.Currency)
	}
	if req.TransactionId != "" {
		db = db.Where("transaction_id = ?", req.TransactionId)
	}
	if req.FromBase58 != "" {
		db = db.Where("from_base58 = ?", req.FromBase58)
	}
	if req.ToBase58 != "" {
		db = db.Where("to_base58 = ?", req.ToBase58)
	}
	if req.Start > 0 {
		db = db.Where("created_at >= ?", req.Start)
	}
	if req.End > 0 {
		db = db.Where("created_at <= ?", req.End)
	}

	items, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("tron trans list find failed, err:%v", err)
		err = biz.TronTransactionFindFailed
		return
	}
	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("tron trans list count failed, err:%v", err)
		err = biz.TronTransactionCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}
