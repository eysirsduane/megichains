package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChainService struct {
	cfg     *global.BackendesConfig
	db      *gorm.DB
	running bool
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
	// case global.ChainNameTron:
	// 	port = s.cfg.Tron.WssNetwork
	default:
		return nil, fmt.Errorf("unknown chain:%v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("evm client Dial failed, chain:%v, err:%v", chain, err)
		return
	}

	return
}

func (s *ChainService) ScanAddressesFunds() {
	if s.running {
		return
	}
	s.running = true
	defer func() {
		s.running = false
	}()

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
			s.tronFunds(addr.Address)
		default:
			logx.Errorf("scan addresses funds found unknown chain, chain:%v", addr.Chain)
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func (s *ChainService) evmFunds(addr string, chain global.ChainName) {
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

	s.updateDB(addr, chain, usdt, usdc)
}

func (s *ChainService) tronFunds(addr string) {
	taddr, err := s.getContractAddress(global.CurrencyTypoUsdt, string(global.ChainNameTron))
	if err != nil {
		logx.Errorf("tron fund get usdc contract address failed, err:%v", err)
		return
	}
	caddr, err := s.getContractAddress(global.CurrencyTypoUsdc, string(global.ChainNameTron))
	if err != nil {
		logx.Errorf("tron fund get usdc contract address failed, err:%v", err)
		return
	}

	usdt, err := s.GetTRC20(addr, taddr)
	if err != nil {
		logx.Errorf("tron fund get usdt balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}
	usdc, err := s.GetTRC20(addr, caddr)
	if err != nil {
		logx.Errorf("tron fund get usdc balance failed, chain:%v, err:%v", global.ChainNameTron, err)
		return
	}

	s.updateDB(addr, global.ChainNameTron, usdt, usdc)
}

func (s *ChainService) updateDB(addr string, chain global.ChainName, usdt, usdc *big.Int) {
	ctx := context.Background()

	db := gorm.G[entity.AddressFund](s.db)
	fund, err := db.Where("address = ?", addr).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		afund := &entity.AddressFund{
			Address: addr,
		}
		switch chain {
		case global.ChainNameBsc:
			afund.BscUsdt = global.Amount(usdt.Int64(), global.AmountTypoBsc)
			afund.BscUsdc = global.Amount(usdc.Int64(), global.AmountTypoBsc)
			afund.Chain = string(global.ChainNameEvm)
		case global.ChainNameEth:
			afund.EthUsdt = global.Amount(usdt.Int64(), global.AmountTypoEth)
			afund.EthUsdc = global.Amount(usdc.Int64(), global.AmountTypoEth)
			afund.Chain = string(global.ChainNameEvm)
		case global.ChainNameTron:
			afund.TronUsdt = global.Amount(usdt.Int64(), global.AmountTypoTron)
			afund.TronUsdc = global.Amount(usdc.Int64(), global.AmountTypoTron)
			afund.Chain = string(global.ChainNameTron)
		}
		err = db.Create(ctx, afund)
		if err != nil {
			logx.Errorf("address fund create failed, err:%v", err)
			return
		}
	} else if err != nil {
		logx.Errorf("address fund first failed, err:%v", err)
	} else {
		switch chain {
		case global.ChainNameBsc:
			fund.BscUsdt = global.Amount(usdt.Int64(), global.AmountTypoBsc)
			fund.BscUsdc = global.Amount(usdc.Int64(), global.AmountTypoBsc)
		case global.ChainNameEth:
			fund.EthUsdt = global.Amount(usdt.Int64(), global.AmountTypoEth)
			fund.EthUsdc = global.Amount(usdc.Int64(), global.AmountTypoEth)
		case global.ChainNameTron:
			fund.TronUsdt = global.Amount(usdt.Int64(), global.AmountTypoTron)
			fund.TronUsdc = global.Amount(usdc.Int64(), global.AmountTypoTron)
		}

		_, err = db.Updates(ctx, fund)
		if err != nil {
			logx.Errorf("address fund first updates failed, err:%v", err)
		}
	}
}

func (s *ChainService) post(path string, body any) (byts []byte, err error) {
	b, _ := json.Marshal(body)
	resp, err := http.Post(s.cfg.Tron.HttpNetwork+path, "application/json", bytes.NewBuffer(b))
	if err != nil {
		logx.Errorf("scan funds tron http post failed, err:%v", err)
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *ChainService) Base58ToHex(addr string) (string, error) {
	decoded := base58.Decode(addr)
	if len(decoded) < 5 {
		return "", errors.New("invalid address length")
	}

	payload := decoded[:len(decoded)-4] // 去掉 checksum
	checksum := decoded[len(decoded)-4:]

	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	if string(checksum) != string(hash[:4]) {
		return "", errors.New("checksum mismatch")
	}

	if payload[0] != 0x41 {
		return "", errors.New("invalid TRON version byte")
	}

	return hex.EncodeToString(payload), nil
}

func (s *ChainService) GetTRC20(addr, contract string) (balance *big.Int, err error) {
	userHex, _ := s.Base58ToHex(addr)
	contractHex, _ := s.Base58ToHex(contract)

	param := "000000000000000000000000" + userHex[2:]

	resp, err := s.post("/wallet/triggersmartcontract", map[string]any{
		"owner_address":     userHex,
		"contract_address":  contractHex,
		"function_selector": "balanceOf(address)",
		"parameter":         param,
	})
	if err != nil {
		return
	}

	var r struct {
		ConstantResult []string `json:"constant_result"`
	}
	json.Unmarshal(resp, &r)

	if len(r.ConstantResult) == 0 {
		return big.NewInt(0), nil
	}

	balance = new(big.Int)
	balance.SetString(r.ConstantResult[0], 16)

	return
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
