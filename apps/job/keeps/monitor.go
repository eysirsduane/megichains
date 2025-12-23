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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	MonitorClientCount            = 100
	MonitorClientSingleQueryLimit = 20
)

type ChainMonitor struct {
	cfg         *global.Config
	clients     sync.Map
	clilen      int
	Receivers   sync.Map
	evmservice  *service.EvmService
	solaservice *service.SolanaService
	addrservice *service.AddressService
}

type ClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *ethclient.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

func NewChainMonitor(cfg *global.Config, evmservice *service.EvmService, addrservice *service.AddressService, solaservice *service.SolanaService) *ChainMonitor {
	monitor := &ChainMonitor{
		cfg:         cfg,
		clients:     sync.Map{},
		Receivers:   sync.Map{},
		evmservice:  evmservice,
		addrservice: addrservice,
		solaservice: solaservice,
	}

	return monitor
}

func (m *ChainMonitor) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
	port := ""
	switch chain {
	case global.ChainNameBsc:
		port = fmt.Sprintf("%v%v", m.cfg.Bsc.GrpcNetwork, m.cfg.Bsc.ApiKey)
	case global.ChainNameEth:
		port = fmt.Sprintf("%v%v", m.cfg.Eth.GrpcNetwork, m.cfg.Eth.ApiKey)
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
	port := fmt.Sprintf("%v%v", m.cfg.Solana.GrpcNetwork, m.cfg.Bsc.ApiKey)
	conn, _, err = websocket.DefaultDialer.Dial(port, nil)
	if err != nil {
		logx.Errorf("Solana chain connect failed, err:%v", err)
		return
	}

	return
}

func (m *ChainMonitor) Listen(chain global.ChainName, oid string, receiver string, seconds int64) {
	m.Receivers.Store(receiver, true)
	defer m.Receivers.Delete(receiver)
	defer m.clearClients()

	switch chain {
	case global.ChainNameEth, global.ChainNameBsc:
		c, err := m.getClientItem(chain)
		if err != nil {
			logx.Errorf("EVM 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*EvmClientItem)
		item.RunningQueryCount++

		m.listenEvm(chain, oid, receiver, seconds, item)
	case global.ChainNameSolana:
		c, err := m.getClientItem(chain)
		if err != nil {
			logx.Errorf("SOLANA 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*SolanaClientItem)
		item.RunningQueryCount++

		m.listenSolana(chain, oid, receiver, seconds, item)
	default:
		logx.Errorf("监听未知的链")
	}

	logx.Infof("chain 事务结束, chain:%v, clen:%v, from:%v", chain, m.clilen, receiver)
}

func (m *ChainMonitor) getClientItem(chain global.ChainName) (item any, err error) {
	found := false
	emu.Lock()
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

		return true
	})
	emu.Unlock()

	if !found {
		for {
			if m.clilen < MonitorClientCount {
				switch chain {
				case global.ChainNameEth, global.ChainNameBsc:
					client, err1 := m.newEvmClient(chain)
					if err1 != nil {
						logx.Errorf("EVM chain Dial 失败, 尝试重连 err:%v", err1)
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

	return
}

func (m *ChainMonitor) listenEvm(chain global.ChainName, oid, receiver string, seconds int64, item *EvmClientItem) {
	contracts := make([]common.Address, 0, 2)
	switch chain {
	case global.ChainNameEth:
		for _, addr := range m.cfg.ContractAddresses {
			if global.ChainName(addr.Chain) == global.ChainNameEth {
				contracts = append(contracts, common.HexToAddress(addr.Address))
			}
		}
	case global.ChainNameBsc:
		for _, addr := range m.cfg.ContractAddresses {
			if global.ChainName(addr.Chain) == global.ChainNameBsc {
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
	go item.listen(ctx, chain, ichan, sub, logs, receiver)

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
	addr := solana.MustPublicKeyFromBase58(
		receiver,
	)

	usdtMintPK := solana.MustPublicKeyFromBase58(m.cfg.Solana.UsdtMint)
	ata, _, err := solana.FindAssociatedTokenAddress(
		addr,
		usdtMintPK,
	)
	if err != nil {
		log.Fatal("find ATA error:", err)
		return
	}

	fmt.Println("📌 Listening USDT ATA:", ata.String())

	// 订阅该 ATA 的账户变化
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "accountSubscribe",
		"params": []interface{}{
			ata.String(),
			map[string]interface{}{
				"encoding":   "jsonParsed",
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
	go item.listen(ctx, ichan, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err = m.saveSolanaOrder(order)
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

	log.Println("✅ Listening USDT transfers for wallet...")
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
		logx.Errorf("EVM chain 记录转账交易失败: err:%v \n", err)
		return
	}

	return
}

func (m *ChainMonitor) saveSolanaOrder(order *entity.SolanaOrder) (err error) {
	err = m.solaservice.Save(order)
	if err != nil {
		logx.Errorf("EVM chain 记录转账交易失败: err:%v \n", err)
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
		req.Seconds = 120

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

		time.Sleep(100 * time.Millisecond)
	}
}
