package keeps

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	case global.ChainNameEth:
	case global.ChainNameBsc:
		item := &EvmClientItem{}
		found := false
		emu.Lock()
		m.clients.Range(func(key, val any) bool {
			cli, ok := val.(*EvmClientItem)
			if ok {
				if cli.Chain == chain && cli.RunningQueryCount < MonitorClientSingleQueryLimit {
					item = cli
					found = true

					return false
				}
			}

			return true
		})
		item.RunningQueryCount++
		emu.Unlock()

		if !found {
			for {
				if m.clilen < MonitorClientCount {
					client, err := m.newEvmClient(chain)
					if err != nil {
						logx.Errorf("EVM chain Dial 失败, 尝试重连 err:%v", err)
						continue
					}
					name := uuid.NewString()
					item.Cfg = m.cfg
					item.Chain = chain
					item.Name = name
					item.Client = client
					m.clients.Store(name, item)
					m.clilen++
					logx.Infof("EVM chain 新增客户端, cname:%v, clen:%v", name, m.clilen)
					break
				} else {
					logx.Errorf("EVM chain 最大客户端已达到上限")
					return
				}
			}
		}

		m.listenEvm(chain, oid, receiver, seconds, item)
	case global.ChainNameSolana:
		item := &SolanaClientItem{}
		found := false
		emu.Lock()
		m.clients.Range(func(key, val any) bool {
			cli, ok := val.(*SolanaClientItem)
			if ok {
				if cli.Chain == chain && cli.RunningQueryCount < MonitorClientSingleQueryLimit {
					item = cli
					found = true

					return false
				}
			}

			return true
		})
		item.RunningQueryCount++
		emu.Unlock()

		if !found {
			for {
				if m.clilen < MonitorClientCount {
					client, err := m.newSolanaClient(chain)
					if err != nil {
						logx.Errorf("SOLANA chain Dial 失败, 尝试重连 err:%v", err)
						continue
					}
					name := uuid.NewString()
					item.Cfg = m.cfg
					item.Chain = chain
					item.Name = name
					item.Client = client
					m.clients.Store(name, item)
					m.clilen++
					logx.Infof("SOLANA chain 新增客户端, cname:%v, clen:%v", name, m.clilen)
					break
				} else {
					logx.Errorf("SOLANA chain 最大客户端已达到上限")
					return
				}
			}
		}

		m.listenSolana(chain, oid, receiver, seconds, item)
	default:
		logx.Errorf("监听未知的链")
	}

	logx.Infof("chain 事务结束, chain:%v, clen:%v, from:%v", chain, m.clilen, receiver)
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

		logx.Infof("🎉🎉🎉 Evm 收到转账, [%v]:[%v], from:%v, to:%v", order.Currency, order.ReceivedAmount, order.FromHex, order.ToHex)

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

		logx.Infof("🎉🎉🎉 Solana 收到转账, [%v]:[%v], from:%v, to:%v", order.Currency, order.ReceivedAmount, order.FromBase58, order.ToBase58)

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
	type ListenRequest struct {
		Chain    string `json:"chain"`
		Receiver string `json:"receiver"`
		Seconds  int64  `json:"seconds"`
	}

	for i := 1; i <= 1000; i++ {
		addr, err := m.addrservice.GetAddress(int64(i))
		if err != nil {
			logx.Errorf("EVM chain 获取监听地址失败, err:%v", err)
			return
		}
		req := &ListenRequest{}
		req.Chain = "EVM chain"
		req.Receiver = addr.AddressHex
		req.Seconds = 1800

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

		time.Sleep(200 * time.Millisecond)
	}
}
