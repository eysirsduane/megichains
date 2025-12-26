package service

import (
	"context"
	"fmt"
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
	db          *gorm.DB
	cfg         *global.Config
	clients     sync.Map
	clilen      int
	qclilen     int
	Receivers   sync.Map
	addrservice *AddressService
	evmservice  *EvmService
	tronservice *TronService
}

func NewChainListenService(cfg *global.Config, db *gorm.DB, addrservice *AddressService, evmservice *EvmService, tronservice *TronService) *ChainListenService {
	return &ChainListenService{
		db:          db,
		cfg:         cfg,
		clients:     sync.Map{},
		Receivers:   sync.Map{},
		addrservice: addrservice,
		evmservice:  evmservice,
		tronservice: tronservice,
	}
}

func (s *ChainListenService) Listen(req *converter.ChainListenReq) {
	key := global.GetOrderAddressKey(string(req.Chain), req.Receiver, req.Currency)
	s.Receivers.Store(key, true)
	defer s.Receivers.Delete(key)
	defer s.clearClients()

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
		logx.Errorf("=== 监听未知的链 ===")
		return
	}

	logx.Infof("req.Chain 事务结束, req.Chain:%v, clen:%v, from:%v", req.Chain, s.clilen, req.Receiver)

	return
}

func (s *ChainListenService) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = fmt.Sprintf("%v%v", s.cfg.Bsc.WssNetwork, s.cfg.Bsc.ApiKey)
	case global.ChainNameEth:
		port = fmt.Sprintf("%v%v", s.cfg.Eth.WssNetwork, s.cfg.Eth.ApiKey)
	default:
		logx.Errorf("未知的链类型: %v", chain)
		return nil, fmt.Errorf("未知的链类型: %v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("EVM req.Chain Dial 失败 err:%v", err)
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
					logx.Errorf("TRON chain Dial 失败, 尝试重连 err:%v", err1)
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
							logx.Errorf("EVM chain Dial 失败, 尝试重连 err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &clients.EvmClientItem{}
						name := uuid.NewString()
						c.Cfg = s.cfg
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						s.clients.Store(name, item)
						s.clilen++
						logx.Infof("EVM chain 新增客户端, cname:%v, clen:%v", name, s.clilen)

						return
					case global.ChainNameSolana:
						client, err1 := s.newSolanaClient(chain)
						if err1 != nil {
							logx.Errorf("SOLANA chain Dial 失败, 尝试重连 err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &clients.SolanaClientItem{}
						name := uuid.NewString()
						c.Cfg = s.cfg
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						s.clients.Store(name, item)
						s.clilen++

						logx.Infof("SOLANA chain 新增客户端, cname:%v, clen:%v", name, s.clilen)

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

func (s *ChainListenService) listenEvm(chain global.ChainName, currency, oid, receiver string, seconds int64, item *clients.EvmClientItem) {
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
		logx.Errorf("未知的链类型: %v", chain)
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
		logx.Errorf("EVM chain 订阅失败:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	ichan := make(chan *entity.EvmOrder, 1)
	go item.Listen(ctx, chain, currency, ichan, sub, logs, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err = s.saveEvmOrder(order)
		if err != nil {
			logx.Errorf("EVM chain 保存日志失败, txid:%v, err:%v", order.TxHash, err)
		}

		logx.Infof("🎉🎉🎉 EVM 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, order.Currency, global.GetFloat64String(order.ReceivedAmount), order.FromHex, order.ToHex)

		err = global.NotifyEPay(s.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TxHash, order.FromHex, order.ToHex, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("EVM chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TxHash, err)
			continue
		}
	}
}

func (s *ChainListenService) listenSolana(chain global.ChainName, oid, receiver string, seconds int64, item *clients.SolanaClientItem) {

}

func (s *ChainListenService) listenTron(chain global.ChainName, currency global.CurrencyTypo, oid, receiver string, seconds int64, item *clients.TronClientItem) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		logx.Errorf("Tron 发现不支持的币种, err:%v", err)
		return
	}

	ichan := make(chan *entity.TronOrder, 1)
	go item.Listen(ctx, ichan, currency, s.cfg.Tron.HttpNetwork, caddr, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err := s.saveTronOrder(order)
		if err != nil {
			logx.Errorf("Tron chain 保存日志失败, txid:%v, err:%v", order.TransactionId, err)
		}

		logx.Infof("🎉🎉🎉 TRON 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, order.Currency, global.GetFloat64String(order.ReceivedAmount), order.FromBase58, order.ToBase58)

		err = global.NotifyEPay(s.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TransactionId, order.FromBase58, order.ToBase58, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("Tron chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TransactionId, err)
			continue
		}
	}
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
		logx.Errorf("Tron 不支持的币种...!, chain:%v, currency:%v", chain, currency)
		err = biz.ContractAddressNotFound
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
				logx.Infof("EVM chain 删除客户端, cname:%v, clen:%v", ecli.Name, s.clilen)
			}
		}

		scli, ok := val.(*clients.SolanaClientItem)
		if ok {
			if scli.RunningQueryCount <= 0 {
				scli.Client.Close()
				s.clients.Delete(key)
				s.clilen--
				logx.Infof("SOLANA chain 删除客户端, cname:%v, clen:%v", scli.Name, s.clilen)
			}

		}

		return true
	})

	emu.Unlock()
}

func (s *ChainListenService) saveEvmOrder(order *entity.EvmOrder) (err error) {
	err = s.evmservice.SaveLog(order)
	if err != nil {
		logx.Errorf("EVM chain 记录转账交易失败: err:%v", err)
		return
	}

	return
}

func (s *ChainListenService) saveTronOrder(order *entity.TronOrder) (err error) {
	err = s.tronservice.SaveOrder(order)
	if err != nil {
		logx.Errorf("TRON chain 记录转账交易失败: err:%v", err)
		return
	}

	return
}

func (s *ChainListenService) saveSolanaOrder(order *entity.SolanaOrder) (err error) {
	// err = m.solaservice.Save(order)
	// if err != nil {
	// 	logx.Errorf("EVM chain 记录转账交易失败: err:%v", err)
	// 	return
	// }

	return
}
