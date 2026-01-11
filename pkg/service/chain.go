package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChainService struct {
	cfg *global.BackendesConfig
	db  *gorm.DB
}

func NewChainService(cfg *global.BackendesConfig, db *gorm.DB) *ChainService {
	return &ChainService{cfg: cfg, db: db}
}

func (s *ChainService) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = s.cfg.Bsc.WssNetwork2
	case global.ChainNameEth:
		port = s.cfg.Eth.WssNetwork2
	default:
		return nil, fmt.Errorf("unknown chain: %v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("EVM client Dial failed, chain:%v, err:%v", chain, err)
		return
	}

	return
}

func (s *ChainService) ScanAddressesFunds() {
	ctx := context.Background()
	addrs, err := gorm.G[entity.Address](s.db).Order("id asc").Find(ctx)
	if err != nil {
		logx.Errorf("scan addresses funds cron address find failed, err:%v", err)
		return
	}

	for _, addr := range addrs {
		switch global.ChainName(addr.Chain) {
		case global.ChainNameEvm:
			s.evmFunds(addr.Address, global.ChainNameBsc)
			s.evmFunds(addr.Address, global.ChainNameEth)
		case global.ChainNameTron:
		default:
			logx.Errorf("scan addresses funds found unknown chain, chain:%v", addr.Chain)
		}

		time.Sleep(time.Millisecond * 200)
	}
}

func (s *ChainService) evmFunds(addr string, chain global.ChainName) {
	ctx := context.Background()
	address := common.HexToAddress(addr)

	taddr, err := s.getContractAddress(global.CurrencyTypoUsdt, string(chain))
	if err != nil {
		logx.Errorf("evm fund get usdt contract address failed, err:%v", err)
		return
	}
	caddr, err := s.getContractAddress(global.CurrencyTypoUsdc, string(chain))
	if err != nil {
		logx.Errorf("evm fund get usdc contract address failed, err:%v", err)
		return
	}

	usdtAddr := common.HexToAddress(taddr)
	usdcAddr := common.HexToAddress(caddr)

	usdt, err := s.ERC20Balance(chain, usdtAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdt balance failed, chain:%v, err:%v", chain, err)
		return
	}
	usdc, err := s.ERC20Balance(chain, usdcAddr, address)
	if err != nil {
		logx.Errorf("evm fund get usdc balance failed, chain:%v, err:%v", chain, err)
		return
	}

	db := gorm.G[entity.AddressFund](s.db)
	fund, err := db.Where("address = ?", addr).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		afund := &entity.AddressFund{
			Address: addr,
			Chain:   string(global.ChainNameEvm),
		}
		switch chain {
		case global.ChainNameBsc:
			afund.BscUsdt = global.Amount(usdt.Int64(), global.AmountTypoBsc)
			afund.BscUsdc = global.Amount(usdc.Int64(), global.AmountTypoBsc)
		case global.ChainNameEth:
			afund.EthUsdt = global.Amount(usdt.Int64(), global.AmountTypoEth)
			afund.EthUsdc = global.Amount(usdc.Int64(), global.AmountTypoEth)
		}
		err = db.Create(ctx, afund)
		if err != nil {
			logx.Errorf("bsc fund create failed, err:%v", err)
			return
		}
	} else if err != nil {
		logx.Errorf("bsc fund first address fund failed, err:%v", err)
	} else {
		switch chain {
		case global.ChainNameBsc:
			fund.BscUsdt = global.Amount(usdt.Int64(), global.AmountTypoBsc)
			fund.BscUsdc = global.Amount(usdc.Int64(), global.AmountTypoBsc)
		case global.ChainNameEth:
			fund.EthUsdt = global.Amount(usdt.Int64(), global.AmountTypoEth)
			fund.EthUsdc = global.Amount(usdc.Int64(), global.AmountTypoEth)
		}
		_, err = db.Updates(ctx, fund)
		if err != nil {
			logx.Errorf("bsc fund first address fund updates failed, err:%v", err)
		}
	}

}

func (s *ChainService) ethFunds(addr string, chain global.ChainName) {
}

func (s *ChainService) getContractAddress(currency global.CurrencyTypo, chain string) (caddr string, err error) {
	switch currency {
	case global.CurrencyTypoUsdt:
		for _, addr := range s.cfg.ContractAddresses {
			if strings.EqualFold(addr.Chain, string(chain)) && strings.EqualFold(addr.Currency, string(currency)) {
				caddr = addr.Address
			}
		}
	case global.CurrencyTypoUsdc:
		for _, addr := range s.cfg.ContractAddresses {
			if strings.EqualFold(addr.Chain, string(chain)) && strings.EqualFold(addr.Currency, string(currency)) {
				caddr = addr.Address
			}
		}
	default:
		err = fmt.Errorf("contract address get failed...!, chain:%v, currency:%v", chain, currency)
		return
	}

	return
}

func (s *ChainService) ERC20Balance(chain global.ChainName, token common.Address, owner common.Address) (*big.Int, error) {
	var cli *ethclient.Client
	switch chain {
	case global.ChainNameBsc:
		cli, _ = s.newEvmClient(global.ChainNameBsc)
	case global.ChainNameEth:
		cli, _ = s.newEvmClient(global.ChainNameEth)
	}
	defer cli.Close()

	erc20ABI := `[{
        "constant":true,
        "inputs":[{"name":"_owner","type":"address"}],
        "name":"balanceOf",
        "outputs":[{"name":"","type":"uint256"}],
        "type":"function"
    }]`

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &token,
		Data: data,
	}

	result, err := cli.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	outputs, err := parsedABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, err
	}

	return outputs[0].(*big.Int), nil
}
