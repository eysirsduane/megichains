package service

import (
	"context"
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

func (s *FundService) Find(ctx context.Context, req *converter.AddressFundListReq) (resp *converter.RespConverter[entity.AddressFund], err error) {
	db := gorm.G[entity.AddressFund](s.db).Order("tron_usdt desc, bsc_usdt desc, eth_usdt desc")
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
	err = gorm.G[entity.AddressFund](s.db).Select("sum(tron_usdt) as tron_usdt, sum(tron_usdc) as tron_usdc, sum(bsc_usdt) as bsc_usdt, sum(bsc_usdc) as bsc_usdc, sum(eth_usdt) as eth_usdt, sum(eth_usdc) as eth_usdc").Scan(ctx, resp)
	if err != nil {
		logx.Errorf("address fund statistics failed, err:%v", err)
	}

	return
}
