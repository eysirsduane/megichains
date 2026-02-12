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

type SolanaService struct {
	db *gorm.DB
}

func NewSolanaService(db *gorm.DB) *SolanaService {
	return &SolanaService{db: db}
}

func (s *SolanaService) Save(order *entity.MerchOrder) (err error) {
	err = s.db.Create(order).Error
	if err != nil {
		logx.Errorf("db solana service save transaction failed, err:%v", err)
		err = biz.SolanaOrderSaveFailed
	}

	return
}

func (s *SolanaService) CreateAddresses(addrs []*entity.Address) (err error) {
	err = s.db.Create(addrs).Error
	if err != nil {
		logx.Errorf("db solana service create address failed, err:%v", err)
		err = biz.AddressCreateFailed
	}
	return
}

func (s *SolanaService) Find(ctx context.Context, req *converter.SolanaTransListReq) (resp *converter.RespConverter[entity.SolanaTransaction], err error) {
	db := gorm.G[entity.SolanaTransaction](s.db).Order("updated_at desc")
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
		logx.Errorf("solana trans list find failed, err:%v", err)
		err = biz.SolanaTransactionFindFailed
		return
	}
	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("solana trans list count failed, err:%v", err)
		err = biz.SolanaTransactionCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}
