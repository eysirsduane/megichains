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

type MerchOrderService struct {
	db *gorm.DB
}

func NewMerchOrderService(db *gorm.DB) *MerchOrderService {
	return &MerchOrderService{db: db}
}

func (s *MerchOrderService) Save(order *entity.MerchantOrder) (err error) {
	return s.db.Save(order).Error
}

func (s *MerchOrderService) LogSave(log *entity.MerchantOrderNotifyLog) (err error) {
	return s.db.Save(log).Error
}

func (s *MerchOrderService) Get(id int64) (order *entity.MerchantOrder, err error) {
	order = &entity.MerchantOrder{}
	err = s.db.Model(&entity.MerchantOrder{}).Where("id = ?", id).First(order).Error
	return
}

func (s *MerchOrderService) GetByMerchId(moid string) (order *entity.MerchantOrder, err error) {
	order = &entity.MerchantOrder{}
	err = s.db.Model(&entity.MerchantOrder{}).Where("merch_order_id = ?", moid).First(order).Error
	return
}

func (s *MerchOrderService) Find(ctx context.Context, req *converter.OrderListReq) (resp *converter.RespConverter[entity.MerchantOrder], err error) {
	db := gorm.G[entity.MerchantOrder](s.db).Order("updated_at desc")
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Typo != "" {
		db = db.Where("typo = ?", req.Typo)
	}
	if req.Mode != "" {
		db = db.Where("mode = ?", req.Mode)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Currency != "" {
		db = db.Where("currency = ?", req.Currency)
	}
	if req.OrderNo != "" {
		db = db.Where("order_no = ?", req.OrderNo)
	}
	if req.MerchantOrderNo != "" {
		db = db.Where("merchant_order_no = ?", req.MerchantOrderNo)
	}
	if req.TransactionId != "" {
		db = db.Where("transaction_id = ?", req.TransactionId)
	}
	if req.FromAddress != "" {
		db = db.Where("from_address = ?", req.FromAddress)
	}
	if req.ToAddress != "" {
		db = db.Where("to_address = ?", req.ToAddress)
	}
	if req.Start > 0 {
		db = db.Where("created_at >= ?", req.Start)
	}
	if req.End > 0 {
		db = db.Where("created_at <= ?", req.End)
	}

	orders, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("order list find failed, err:%v", err)
		err = biz.OrderFindFailed
		return
	}
	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("order list count failed, err:%v", err)
		err = biz.OrderCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(orders, req.Current, req.Size, total)

	return
}
