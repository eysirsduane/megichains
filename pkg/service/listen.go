package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service/clients"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var emu sync.Mutex
var tmu sync.Mutex

const (
	MonitorClientCount             = 100
	MonitorClientSingleQueryLimit  = 20
	BitQueryClientSingleQueryLimit = 20
)

type ChainListenService struct {
	db           *gorm.DB
	cfg          *global.BackendesConfig
	clients      sync.Map
	clilen       int
	qclilen      int
	Receivers    sync.Map
	addrservice  *AddressService
	orderservice *MerchOrderService
	evmservice   *EvmService
	tronservice  *TronService
}

func NewChainListenService(cfg *global.BackendesConfig, db *gorm.DB, addrservice *AddressService, orderservice *MerchOrderService, evmservice *EvmService, tronservice *TronService) *ChainListenService {
	return &ChainListenService{
		db:           db,
		cfg:          cfg,
		clients:      sync.Map{},
		Receivers:    sync.Map{},
		addrservice:  addrservice,
		orderservice: orderservice,
		evmservice:   evmservice,
		tronservice:  tronservice,
	}
}

func (s *ChainListenService) Listen(req *converter.ChainListenReq) {
	key := global.GetOrderAddressKey(string(req.Chain), req.Receiver, req.Currency)
	s.Receivers.Store(key, true)
	defer s.Receivers.Delete(key)
	defer s.clearClients()

	req.Seconds += 120
	switch req.Chain {
	case global.ChainNameEth, global.ChainNameBsc:
		emu.Lock()
		c, err1 := s.getClientItem(req.Chain)
		if err1 != nil {
			logx.Errorf("EVM 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*clients.EvmClientItem)
		item.RunningQueryCount++
		emu.Unlock()

		s.listenEvm(req.Chain, req.Currency, req.MerchOrderId, req.Receiver, req.Seconds, item)
	case global.ChainNameSolana:
		emu.Lock()
		c, err1 := s.getClientItem(req.Chain)
		if err1 != nil {
			logx.Errorf("SOLANA 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*clients.SolanaClientItem)
		item.RunningQueryCount++
		emu.Unlock()

		s.listenSolana(req.Chain, req.MerchOrderId, req.Receiver, req.Seconds, item)
	case global.ChainNameTron:
		tmu.Lock()
		c, err1 := s.getClientItem(req.Chain)
		if err1 != nil {
			logx.Errorf("TRON 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*clients.TronClientItem)
		tmu.Unlock()

		s.listenTron(req.Chain, global.CurrencyTypo(req.Currency), req.MerchOrderId, req.Receiver, req.Seconds, item)
	default:
		logx.Errorf("=== 监听了未知的链 ===")
		return
	}

	logx.Infof("监听事务结束, chain:%v, clen:%v, from:%v", req.Chain, s.clilen, req.Receiver)
}

func (s *ChainListenService) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = fmt.Sprintf("%v%v", s.cfg.Bsc.WssNetwork, s.cfg.Bsc.ApiKey)
	case global.ChainNameEth:
		port = fmt.Sprintf("%v%v", s.cfg.Eth.WssNetwork, s.cfg.Eth.ApiKey)
	default:
		return nil, fmt.Errorf("unknown chain: %v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("EVM req.Chain Dial failed err:%v", err)
		return
	}

	return
}

func (s *ChainListenService) newSolanaClient(chain global.ChainName) (conn *websocket.Conn, err error) {
	port := fmt.Sprintf("%v%v", s.cfg.Solana.WssNetwork, s.cfg.Solana.ApiKey)
	conn, _, err = websocket.DefaultDialer.Dial(port, nil)
	if err != nil {
		logx.Errorf("Solana req.Chain connect failed, err:%v", err)
		return
	}

	return
}

func (s *ChainListenService) newTronClient(chain global.ChainName) (conn *websocket.Conn, err error) {
	// port := fmt.Sprintf("%v", s.cfg.Tron.HttpNetwork)
	// conn, _, err = websocket.DefaultDialer.Dial(port, nil)
	// if err != nil {
	// 	logx.Errorf("Tron req.Chain connect failed, err:%v", err)
	// 	return
	// }

	return
}

func (s *ChainListenService) getClientItem(chain global.ChainName) (item any, err error) {
	found := false
	s.clients.Range(func(key, val any) bool {
		ecli, ok := val.(*clients.EvmClientItem)
		if ok {
			if ecli.Chain == chain && ecli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = ecli
				found = true

				return false
			}
		}

		scli, ok := val.(*clients.SolanaClientItem)
		if ok {
			if scli.Chain == chain && scli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = scli
				found = true

				return false
			}
		}

		tcli, ok := val.(*clients.TronClientItem)
		if ok {
			if tcli.Chain == chain && tcli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = tcli
				found = true

				return false
			}
		}

		return true
	})

	if !found {
		for {
			if chain == global.ChainNameTron {
				client, err1 := s.newTronClient(chain)
				if err1 != nil {
					logx.Errorf("TRON chain Dial failed, try again err:%v", err1)
					time.Sleep(time.Second * 2)
					continue
				}

				c := &clients.TronClientItem{}
				name := uuid.NewString()
				c.Chain = chain
				c.Name = name
				c.Client = client
				item = c

				logx.Infof("TRON chain 新增客户端, cname:%v, clen:%v", name, s.clilen)

				return
			} else {
				if s.clilen < MonitorClientCount {
					switch chain {
					case global.ChainNameEth, global.ChainNameBsc:
						client, err1 := s.newEvmClient(chain)
						if err1 != nil {
							logx.Errorf("EVM chain Dial failed, try again err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &clients.EvmClientItem{}
						name := uuid.NewString()
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						s.clients.Store(name, item)
						s.clilen++
						logx.Infof("EVM chain 创建新客户端, cname:%v, clen:%v", name, s.clilen)

						return
					case global.ChainNameSolana:
						client, err1 := s.newSolanaClient(chain)
						if err1 != nil {
							logx.Errorf("SOLANA chain Dial failed, try again err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &clients.SolanaClientItem{}
						name := uuid.NewString()
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						s.clients.Store(name, item)
						s.clilen++

						logx.Infof("SOLANA chain 创建新客户端, cname:%v, clen:%v", name, s.clilen)

						return
					}
				} else {
					err = biz.ChainClientUpToMaxCount
					logx.Errorf("最大客户端已达到上限, chain:%v", chain)
					return
				}
			}
		}
	}

	return
}

func (s *ChainListenService) getContractAddress(currency global.CurrencyTypo, chain string) (caddr string, err error) {
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

func (s *ChainListenService) clearClients() {
	emu.Lock()

	s.clients.Range(func(key, val any) bool {
		ecli, ok := val.(*clients.EvmClientItem)
		if ok {
			if ecli.RunningQueryCount <= 0 {
				ecli.Client.Close()
				s.clients.Delete(key)
				s.clilen--
				logx.Infof("EVM chain delete client, cname:%v, clen:%v", ecli.Name, s.clilen)
			}
		}

		scli, ok := val.(*clients.SolanaClientItem)
		if ok {
			if scli.RunningQueryCount <= 0 {
				scli.Client.Close()
				s.clients.Delete(key)
				s.clilen--
				logx.Infof("SOLANA chain delete client, cname:%v, clen:%v", scli.Name, s.clilen)
			}

		}

		return true
	})

	emu.Unlock()
}

func (s *ChainListenService) listenEvm(chain global.ChainName, currency, moid, receiver string, seconds int64, item *clients.EvmClientItem) {
	contracts := make([]common.Address, 0, 1)
	switch chain {
	case global.ChainNameEth:
		for _, addr := range s.cfg.ContractAddresses {
			if global.ChainName(addr.Chain) == global.ChainNameEth && addr.Currency == currency {
				contracts = append(contracts, common.HexToAddress(addr.Address))
			}
		}
	case global.ChainNameBsc:
		for _, addr := range s.cfg.ContractAddresses {
			if global.ChainName(addr.Chain) == global.ChainNameBsc && addr.Currency == currency {
				contracts = append(contracts, common.HexToAddress(addr.Address))
			}
		}
	default:
		logx.Errorf("Evm contract address get failed, chain:%v, currency:%v", chain, currency)
		return
	}

	transferSig := crypto.Keccak256Hash(
		[]byte("Transfer(address,address,uint256)"),
	)

	query := ethereum.FilterQuery{
		Addresses: contracts,
		Topics: [][]common.Hash{
			{transferSig},
			nil,
			{common.HexToHash(receiver)},
		},
	}

	logs := make(chan types.Log)
	sub, err := item.Client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logx.Errorf("EVM chain subscribe failed:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	order := &entity.MerchOrder{
		MerchOrderId: moid,
		Chain:        string(chain),
		Typo:         string(global.OrderTypoIn),
		Status:       string(global.OrderStatusCreated),
		Currency:     currency,
		ToAddress:    receiver,
		Description:  "",
	}
	err = s.orderservice.Save(order)
	if err != nil {
		logx.Errorf("order service evm order create failed, moid:%v, receiver:%v, err:%v", moid, receiver, err)
		return
	}

	ichan := make(chan *entity.EvmLog, 1)
	go item.Listen(ctx, chain, ichan, currency, sub, logs, receiver)

	for log := range ichan {
		log.Chain = string(chain)
		err = s.evmservice.LogSave(log)
		if err != nil {
			logx.Errorf("evm service save log failed, err:%v", err)
			err = biz.EvmOrderSaveFailed
			continue
		}

		logx.Infof("🎉🎉🎉 EVM 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, log.Currency, global.GetFloat64String(log.Amount), log.FromHex, log.ToHex)

		order.LogId = log.Id
		order.MerchOrderId = moid
		order.TransactionId = log.TxHash
		order.FromAddress = log.FromHex
		order.ReceivedAmount = log.Amount
		order.ReceivedSun = log.Sun
		order.Status = string(global.OrderStatusSuccess)

		err = global.NotifyEPay(s.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TransactionId, order.FromAddress, order.ToAddress, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("EVM chain notify failed, moid:%v, txid:%v, err:%v", moid, log.TxHash, err)
			order.Status = string(global.OrderStatusNotifyFailed)
			order.Description = err.Error()
		}

		err = s.orderservice.Save(order)
		if err != nil {
			logx.Errorf("EVM chain order service save order failed, moid:%v, txid:%v, err:%v", moid, order.TransactionId, err)
			return
		}
	}
}

func (s *ChainListenService) listenSolana(chain global.ChainName, oid, receiver string, seconds int64, item *clients.SolanaClientItem) {

}

func (s *ChainListenService) listenTron(chain global.ChainName, currency global.CurrencyTypo, moid, receiver string, seconds int64, item *clients.TronClientItem) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		logx.Errorf("Tron contract address get failed, chain:%v, currency:%v, err:%v", chain, currency, err)
		err = biz.ContractAddressNotFound
		return
	}

	order := &entity.MerchOrder{
		MerchOrderId: moid,
		Chain:        string(chain),
		Typo:         string(global.OrderTypoIn),
		Status:       string(global.OrderStatusCreated),
		Currency:     string(currency),
		ToAddress:    receiver,
		Description:  "",
	}
	err = s.orderservice.Save(order)
	if err != nil {
		logx.Errorf("order service create tron order failed, moid:%v, receiver:%v, err:%v", moid, receiver, err)
		err = biz.TronOrderSaveFailed
		return
	}

	ichan := make(chan *entity.TronTransaction, 1)
	go item.Listen(ctx, chain, ichan, currency, s.cfg.Tron.HttpNetwork, caddr, receiver)

	for trans := range ichan {
		trans.Chain = string(chain)
		err = s.tronservice.TransSave(trans)
		if err != nil {
			logx.Errorf("tron service save trans failed, err:%v", err)
			err = biz.TronTransSaveFailed
			continue
		}

		logx.Infof("🎉🎉🎉 TRON 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, trans.Currency, global.GetFloat64String(trans.Amount), trans.FromBase58, trans.ToBase58)

		order.LogId = trans.Id
		order.MerchOrderId = moid
		order.TransactionId = trans.TransactionId
		order.Chain = string(chain)
		order.FromAddress = trans.FromBase58
		order.ReceivedAmount = trans.Amount
		order.ReceivedSun = trans.Sun
		order.Status = string(global.OrderStatusSuccess)

		err = global.NotifyEPay(s.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TransactionId, order.FromAddress, order.ToAddress, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("Tron chain notify failed, moid:%v, txid:%v, err:%v", moid, trans.TransactionId, err)
			order.Status = string(global.OrderStatusNotifyFailed)
			order.Description = err.Error()
		}

		err = s.orderservice.Save(order)
		if err != nil {
			logx.Errorf("Tron chain order service save order failed, moid:%v, txid:%v, err:%v", moid, order.TransactionId, err)
			err = biz.TronOrderSaveFailed
			return
		}
	}
}

func (s *ChainListenService) ListenMany() {
	cnames := []string{"BSC", "ETH", "TRON"}
	currencys := []string{"USDT", "USDC"}

	for i := 1; i < 500; i++ {
		chain := ""
		iaddr := rand.IntN(1000)
		addr, _ := s.addrservice.GetAddress(int64(iaddr))
		switch addr.Chain {
		case "EVM":
			iname := rand.IntN(2)
			chain = cnames[iname]
		case "TRON":
			chain = cnames[2]
		default:
			continue
		}

		icurr := rand.IntN(2)
		currency := currencys[icurr]

		go s.Listen(&converter.ChainListenReq{
			MerchOrderId: uuid.NewString(),
			Chain:        global.ChainName(chain),
			Currency:     currency,
			Receiver:     addr.Address,
			Seconds:      180,
		})

		time.Sleep(time.Millisecond * 100)
	}
}
