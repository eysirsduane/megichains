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

type FundService struct {
	db *gorm.DB
}

func NewFundService(db *gorm.DB) *FundService {
	return &FundService{db: db}
}

func (s *FundService) Find(ctx context.Context, req *converter.AddressFundListReq) (resp *converter.RespConverter[entity.Address], err error) {
	db := gorm.G[entity.Address](s.db).Order("tron_usdt desc, bsc_usdt desc, eth_usdt desc")
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Address != "" {
		db = db.Where("address = ?", req.Address)
	}
	funds, err := db.Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Find(ctx)
	if err != nil {
		logx.Errorf("fund find failed, err:%v", err)
		return
	}

	total, err := db.Count(ctx, "id")
	if err != nil {
		logx.Errorf("fund count failed, err:%v", err)
		return
	}

	resp = converter.ConvertToPagingRecordsResp(funds, req.Current, req.Size, total)

	return
}

func (s *FundService) Statistics(ctx context.Context) (resp *converter.AddressFundStatisticsResp, err error) {
	resp = &converter.AddressFundStatisticsResp{}
	err = gorm.G[entity.AddressBalance](s.db).Select("sum(tron_usdt) as tron_usdt, sum(tron_usdc) as tron_usdc, sum(bsc_usdt) as bsc_usdt, sum(bsc_usdc) as bsc_usdc, sum(eth_usdt) as eth_usdt, sum(eth_usdc) as eth_usdc, sum(solana_usdt) as solana_usdt, sum(solana_usdc) as solana_usdc").Scan(ctx, resp)
	if err != nil {
		logx.Errorf("address fund statistics failed, err:%v", err)
	}

	return
}

func (s *FundService) FindCollectList(ctx context.Context, req *converter.AddressFundCollectListReq) (resp *converter.RespConverter[*converter.AddressFundCollectItem], err error) {
	db := s.db.Model(entity.AddressFundCollect{}).Order("id desc")
	if req.ReceiverAddress != "" {
		db = db.Where("receiver_address = ?", req.ReceiverAddress)
	}
	if req.AddressGroupId != 0 {
		db = db.Where("address_group_id = ?", req.AddressGroupId)
	}
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Currency != "" {
		db = db.Where("currency = ?", req.Currency)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	items := make([]*converter.AddressFundCollectItem, 0)
	err = db.Session(&gorm.Session{}).Select("address_fund_collects.*, address_groups.name as address_group_name, users.username as username").Joins("left join address_groups on address_fund_collects.address_group_id = address_groups.id").Joins("left join users on address_fund_collects.user_id = users.id").Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Scan(&items).Error
	if err != nil {
		logx.Errorf("address fund collect list paging failed, err:%v", err)
		err = biz.AddressFundCollectFindFailed
		return
	}
	total := int64(0)
	err = db.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		logx.Errorf("address fund collect list count failed, err:%v", err)
		err = biz.AddressFundCollectCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)
	return
}

func (s *FundService) ScanFundCollectsStatus() {
	collects := make([]*entity.AddressFundCollect, 0)
	err := s.db.Model(&entity.AddressFundCollect{}).Where("status = ?", global.CollectStatusProcessing).Find(&collects).Error
	if err != nil {
		logx.Errorf("job scan fund collect status failed, err:%v", err)
		return
	}

	for _, coll := range collects {
		created := int64(0)
		err = s.db.Model(&entity.AddressFundCollectLog{}).Where("collect_id = ? and status = ?", coll.Id, global.CollectLogStatusCreated).Count(&created).Error
		if err != nil {
			logx.Errorf("job scan fund collect status count successful failed, err:%v", err)
			continue
		}
		if created > 0 {
			continue
		}

		calc := &global.CollectCalc{}
		err = s.db.Model(&entity.AddressFundCollectLog{}).Where("collect_id = ? and status = ?", coll.Id, global.CollectLogStatusSuccess).Select("count(id) as success_count, sum(amount) as success_amount, sum(total_gas_fee) as total_gas_fee").Scan(calc).Error
		if err != nil {
			logx.Errorf("job scan fund collect status count successful failed, err:%v", err)
			continue
		}
		coll.TotalGasFee = calc.TotalGasFee
		coll.TotalGasFeeCurrency = global.Amount(calc.TotalGasFee, global.AmountTypo18e)
		if coll.TotalCount == calc.SuccessCount {
			coll.Status = string(global.CollectLogStatusSuccess)
		} else {
			faileds := int64(0)
			err = s.db.Model(&entity.AddressFundCollectLog{}).Where("collect_id = ? and status = ?", coll.Id, global.CollectLogStatusFailed).Count(&faileds).Error
			if err != nil {
				logx.Errorf("job scan fund collect status count successful failed, err:%v", err)
				continue
			}
			if coll.TotalCount == faileds {
				coll.Status = string(global.CollectStatusFailed)
			} else {
				coll.Status = string(global.CollectStatusPartially)
			}
		}

		coll.SuccessCount = calc.SuccessCount
		coll.SuccessAmount = calc.SuccessAmount
		err = s.db.Save(coll).Error
		if err != nil {
			logx.Errorf("job scan fund collect status save collect failed, err:%v", err)
			continue
		}

		logx.Infof("job scan fund collect handle successful, collect:%v", coll.Id)
	}
}

func (s *FundService) FindCollectLogList(ctx context.Context, req *converter.AddressFundCollectLogListReq) (resp *converter.RespConverter[*converter.AddressFundCollectLogItem], err error) {
	db := s.db.Model(entity.AddressFundCollectLog{}).Order("id desc")
	if req.CollectId > 0 {
		db = db.Where("collect_id = ?", req.CollectId)
	}
	if req.Chain != "" {
		db = db.Where("chain = ?", req.Chain)
	}
	if req.Currency != "" {
		db = db.Where("currency = ?", req.Currency)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.FromAddress != "" {
		db = db.Where("from_address = ?", req.FromAddress)
	}
	if req.ReceiverAddress != "" {
		db = db.Where("receiver_address = ?", req.ReceiverAddress)
	}

	items := make([]*converter.AddressFundCollectLogItem, 0)
	err = db.Session(&gorm.Session{}).Offset(global.Offset(req.Current, req.Size)).Limit(req.Size).Scan(&items).Error
	if err != nil {
		logx.Errorf("address fund collect log paging failed, err:%v", err)
		err = biz.AddressFundCollectLogFindFailed
		return
	}
	total := int64(0)
	err = db.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		logx.Errorf("address fund collect log count failed, err:%v", err)
		err = biz.AddressFundCollectLogCountFailed
		return
	}

	resp = converter.ConvertToPagingRecordsResp(items, req.Current, req.Size, total)

	return
}
