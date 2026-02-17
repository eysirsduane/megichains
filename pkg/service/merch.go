package service

import (
	"context"
	"fmt"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MerchService struct {
	db *gorm.DB
}

func NewMerchService(db *gorm.DB) *MerchService {
	return &MerchService{db: db}
}

func (s *MerchService) Get(id int64) (merch *entity.Merchant, err error) {
	merch = &entity.Merchant{}
	err = s.db.First(merch, id).Error

	return
}

func (s *MerchService) GetByAccount(account string) (merch *entity.Merchant, err error) {
	merch = &entity.Merchant{}
	err = s.db.Model(&entity.Merchant{}).Where("merchant_account = ?", account).First(merch).Error

	return
}

func (s *MerchService) Save(ctx context.Context, req *converter.MerchantItem) (err error) {
	merch := &entity.Merchant{}
	copier.Copy(&merch, req)

	if req.Id > 0 {
		_, err = gorm.G[entity.Merchant](s.db).Updates(ctx, *merch)
		if err != nil {
			logx.Errorf("merchant update failed, id:%v, err:%v", req.Id, err)
			err = biz.MerchantSaveFailed
			return
		}
	} else {
		merch.SecretKey = global.GenerateRandomString()
		merch.MerchantAccount = fmt.Sprintf("M%v", time.Now().UnixMicro())
		err = gorm.G[entity.Merchant](s.db).Create(ctx, merch)
		if err != nil {
			logx.Errorf("merchant create failed, id:%v, err:%v", req.Id, err)
			err = biz.MerchantCreateFailed
			return
		}
	}

	return
}

func (s *MerchService) Find(ctx context.Context, req *converter.MerchantListReq) (resp *converter.RespConverter[*entity.Merchant], err error) {
	db := gorm.G[*entity.Merchant](s.db).Order("id desc")
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if req.MerchantAccount != "" {
		db = db.Where("merchant_account = ?", req.MerchantAccount)
	}

	items, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("merchant list find failed, err:%v", err)
		err = biz.TronTransactionFindFailed
		return
	}
	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("merchant list count failed, err:%v", err)
		err = biz.TronTransactionCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}
