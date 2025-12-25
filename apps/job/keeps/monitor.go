package keeps

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	MonitorClientCount             = 100
	MonitorClientSingleQueryLimit  = 20
	BitQueryClientSingleQueryLimit = 20
)

var emu sync.Mutex
var tmu sync.Mutex

type ChainMonitor struct {
	cfg         *global.Config
	clients     sync.Map
	clilen      int
	qclilen     int
	Receivers   sync.Map
	evmservice  *service.EvmService
	solaservice *service.SolanaService
	addrservice *service.AddressService
	tronservice *service.TronService
}

type ClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *ethclient.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

func NewChainMonitor(cfg *global.Config, evmservice *service.EvmService, addrservice *service.AddressService, solaservice *service.SolanaService, tronservice *service.TronService) *ChainMonitor {
	monitor := &ChainMonitor{
		cfg:         cfg,
		clients:     sync.Map{},
		Receivers:   sync.Map{},
		evmservice:  evmservice,
		addrservice: addrservice,
		solaservice: solaservice,
		tronservice: tronservice,
	}

	return monitor
}

func (m *ChainMonitor) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = fmt.Sprintf("%v%v", m.cfg.Bsc.WssNetwork, m.cfg.Bsc.ApiKey)
	case global.ChainNameEth:
		port = fmt.Sprintf("%v%v", m.cfg.Eth.WssNetwork, m.cfg.Eth.ApiKey)
	default:
		logx.Errorf("未知的链类型: %v", chain)
		return nil, fmt.Errorf("未知的链类型: %v", chain)
	}

	client, err = ethclient.Dial(port)
	if err != nil {
		logx.Errorf("EVM chain Dial 失败 err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) newSolanaClient(chain global.ChainName) (conn *websocket.Conn, err error) {
	port := fmt.Sprintf("%v%v", m.cfg.Solana.WssNetwork, m.cfg.Solana.ApiKey)
	conn, _, err = websocket.DefaultDialer.Dial(port, nil)
	if err != nil {
		logx.Errorf("Solana chain connect failed, err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) newTronClient(chain global.ChainName) (conn *websocket.Conn, err error) {
	// port := fmt.Sprintf("%v", m.cfg.Tron.HttpNetwork)
	// conn, _, err = websocket.DefaultDialer.Dial(port, nil)
	// if err != nil {
	// 	logx.Errorf("Tron chain connect failed, err:%v", err)
	// 	return
	// }

	return
}

func (m *ChainMonitor) Listen(chain global.ChainName, currency, oid string, receiver string, seconds int64) {
	key := global.GetOrderAddressKey(string(chain), receiver, currency)
	m.Receivers.Store(key, true)
	defer m.Receivers.Delete(key)
	defer m.clearClients()

	switch chain {
	case global.ChainNameEth, global.ChainNameBsc:
		emu.Lock()
		c, err := m.getClientItem(chain)
		if err != nil {
			logx.Errorf("EVM 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*EvmClientItem)
		item.RunningQueryCount++
		emu.Unlock()

		m.listenEvm(chain, currency, oid, receiver, seconds, item)
	case global.ChainNameSolana:
		emu.Lock()
		c, err := m.getClientItem(chain)
		if err != nil {
			logx.Errorf("SOLANA 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*SolanaClientItem)
		item.RunningQueryCount++
		emu.Unlock()

		m.listenSolana(chain, oid, receiver, seconds, item)
	case global.ChainNameTron:
		tmu.Lock()
		c, err := m.getClientItem(chain)
		if err != nil {
			logx.Errorf("TRON 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*TronClientItem)
		tmu.Unlock()

		m.listenTron(chain, global.CurrencyTypo(currency), oid, receiver, seconds, item)
	default:
		logx.Errorf("=== 监听未知的链 ===")
		return
	}

	logx.Infof("chain 事务结束, chain:%v, clen:%v, from:%v", chain, m.clilen, receiver)
}

func (m *ChainMonitor) getClientItem(chain global.ChainName) (item any, err error) {
	found := false
	m.clients.Range(func(key, val any) bool {
		ecli, ok := val.(*EvmClientItem)
		if ok {
			if ecli.Chain == chain && ecli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = ecli
				found = true

				return false
			}
		}

		scli, ok := val.(*SolanaClientItem)
		if ok {
			if scli.Chain == chain && scli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = scli
				found = true

				return false
			}
		}

		tcli, ok := val.(*TronClientItem)
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
				client, err1 := m.newTronClient(chain)
				if err1 != nil {
					logx.Errorf("TRON chain Dial 失败, 尝试重连 err:%v", err1)
					time.Sleep(time.Second * 2)
					continue
				}

				c := &TronClientItem{}
				name := uuid.NewString()
				c.Chain = chain
				c.Name = name
				c.Client = client
				item = c

				logx.Infof("TRON chain 新增客户端, cname:%v, clen:%v", name, m.clilen)

				return
			} else {
				if m.clilen < MonitorClientCount {
					switch chain {
					case global.ChainNameEth, global.ChainNameBsc:
						client, err1 := m.newEvmClient(chain)
						if err1 != nil {
							logx.Errorf("EVM chain Dial 失败, 尝试重连 err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &EvmClientItem{}
						name := uuid.NewString()
						c.Cfg = m.cfg
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						m.clients.Store(name, item)
						m.clilen++
						logx.Infof("EVM chain 新增客户端, cname:%v, clen:%v", name, m.clilen)

						return
					case global.ChainNameSolana:
						client, err1 := m.newSolanaClient(chain)
						if err1 != nil {
							logx.Errorf("SOLANA chain Dial 失败, 尝试重连 err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &SolanaClientItem{}
						name := uuid.NewString()
						c.Cfg = m.cfg
						c.Chain = chain
						c.Name = name
						c.Client = client
						item = c
						m.clients.Store(name, item)
						m.clilen++

						logx.Infof("SOLANA chain 新增客户端, cname:%v, clen:%v", name, m.clilen)

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

func (m *ChainMonitor) listenEvm(chain global.ChainName, currency, oid, receiver string, seconds int64, item *EvmClientItem) {
	contracts := make([]common.Address, 0, 1)
	switch chain {
	case global.ChainNameEth:
		for _, addr := range m.cfg.ContractAddresses {
			if global.ChainName(addr.Chain) == global.ChainNameEth && addr.Currency == currency {
				contracts = append(contracts, common.HexToAddress(addr.Address))
			}
		}
	case global.ChainNameBsc:
		for _, addr := range m.cfg.ContractAddresses {
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
	go item.listen(ctx, chain, currency, ichan, sub, logs, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err = m.saveEvmOrder(order)
		if err != nil {
			logx.Errorf("EVM chain 保存日志失败, txid:%v, err:%v", order.TxHash, err)
		}

		logx.Infof("🎉🎉🎉 EVM 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, order.Currency, global.GetFloat64String(order.ReceivedAmount), order.FromHex, order.ToHex)

		err = global.NotifyEPay(m.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TxHash, order.FromHex, order.ToHex, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("EVM chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TxHash, err)
			continue
		}
	}
}

func (m *ChainMonitor) listenSolana(chain global.ChainName, oid, receiver string, seconds int64, item *SolanaClientItem) {
	// 订阅该 ATA 的账户变化
	req := LogsSubscribeReq{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "logsSubscribe",
		Params: []any{
			map[string]any{
				"mentions": []string{receiver},
			},
			map[string]any{
				"commitment": "confirmed",
			},
		},
	}

	if err := item.Client.WriteJSON(req); err != nil {
		log.Fatal("subscribe error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	ichan := make(chan *entity.SolanaOrder, 1)
	go item.listen(ctx, m.cfg.Solana.GrpcNetwork, m.cfg.Solana.UsdtMint, ichan, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err := m.saveSolanaOrder(order)
		if err != nil {
			logx.Errorf("Solana chain 保存日志失败, txid:%v, err:%v", order.TxHash, err)
		}

		logx.Infof("🎉🎉🎉 SOLANA 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, order.Currency, global.GetFloat64String(order.ReceivedAmount), order.FromBase58, order.ToBase58)

		err = global.NotifyEPay(m.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TxHash, order.FromBase58, order.ToBase58, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("Solana chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TxHash, err)
			continue
		}
	}
}

func (m *ChainMonitor) listenTron(chain global.ChainName, currency global.CurrencyTypo, oid, receiver string, seconds int64, item *TronClientItem) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	caddr := ""
	switch currency {
	case global.CurrencyTypoUsdt:
		for _, addr := range m.cfg.ContractAddresses {
			if strings.EqualFold(addr.Chain, string(chain)) && strings.EqualFold(addr.Currency, string(currency)) {
				caddr = addr.Address
			}
		}
	case global.CurrencyTypoUsdc:
	default:
		logx.Errorf("Tron 未知的币种...!, chain:%v, currency:%v", chain, currency)
		return
	}

	ichan := make(chan *entity.TronOrder, 1)
	go item.listen(ctx, ichan, currency, m.cfg.Tron.HttpNetwork, caddr, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err := m.saveTronOrder(order)
		if err != nil {
			logx.Errorf("Tron chain 保存日志失败, txid:%v, err:%v", order.TransactionId, err)
		}

		logx.Infof("🎉🎉🎉 TRON 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, order.Currency, global.GetFloat64String(order.ReceivedAmount), order.FromBase58, order.ToBase58)

		err = global.NotifyEPay(m.cfg.EPay.NotifyUrl, order.MerchOrderId, order.TransactionId, order.FromBase58, order.ToBase58, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("Tron chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TransactionId, err)
			continue
		}
	}
}

func (m *ChainMonitor) GenerateETHAddress() {
	addrs := make([]*entity.Address, 0, 1000)

	for i := 0; i < 1000; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert public key type")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		addr := &entity.Address{
			Chain:       "EVM chain",
			Typo:        "KEEP",
			Status:      "ACTIVE",
			AddressHex:  address,
			Secrect:     privateKeyHex,
			Description: fmt.Sprintf("EVM chain KEEP 地址 %d", i+1),
		}
		addrs = append(addrs, addr)
	}

	m.evmservice.CreateAddresses(addrs)
}

func (m *ChainMonitor) clearClients() {
	emu.Lock()

	m.clients.Range(func(key, val any) bool {
		ecli, ok := val.(*EvmClientItem)
		if ok {
			if ecli.RunningQueryCount <= 0 {
				ecli.Client.Close()
				m.clients.Delete(key)
				m.clilen--
				logx.Infof("EVM chain 删除客户端, cname:%v, clen:%v", ecli.Name, m.clilen)
			}
		}

		scli, ok := val.(*SolanaClientItem)
		if ok {
			if scli.RunningQueryCount <= 0 {
				scli.Client.Close()
				m.clients.Delete(key)
				m.clilen--
				logx.Infof("SOLANA chain 删除客户端, cname:%v, clen:%v", scli.Name, m.clilen)
			}

		}

		return true
	})

	emu.Unlock()
}

func (m *ChainMonitor) saveEvmOrder(order *entity.EvmOrder) (err error) {
	err = m.evmservice.SaveLog(order)
	if err != nil {
		logx.Errorf("EVM chain 记录转账交易失败: err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) saveTronOrder(order *entity.TronOrder) (err error) {
	err = m.tronservice.SaveOrder(order)
	if err != nil {
		logx.Errorf("TRON chain 记录转账交易失败: err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) saveSolanaOrder(order *entity.SolanaOrder) (err error) {
	err = m.solaservice.Save(order)
	if err != nil {
		logx.Errorf("EVM chain 记录转账交易失败: err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) RangeListen() {

	for i := 1; i <= 100; i++ {
		addr, err := m.addrservice.GetAddress(int64(i))
		if err != nil {
			logx.Errorf("EVM chain 获取监听地址失败, err:%v", err)
			return
		}
		req := &global.ListenReq{}
		req.MerchOrderId = "123"
		req.Chain = "ETH"
		req.Receiver = addr.AddressHex
		req.Seconds = 60

		bs, err := json.Marshal(req)
		if err != nil {
			logx.Errorf("EVM chain 监听请求序列化失败, err:%v", err)
			return
		}

		resp, err := http.Post("http://127.0.0.1:7002/listen", "application/json", bytes.NewReader(bs))
		if err != nil {
			logx.Errorf("EVM chain 发送监听请求失败, err:%v", err)
			return
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logx.Errorf("EVM chain 读取监听响应失败, err:%v", err)
			return
		}

		logx.Infof("EVM chain 发送监听请求成功, resp:%s", string(respBytes))

	}
}
