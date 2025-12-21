package keeps

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var emu sync.Mutex

type ETHMonitor struct {
	cfg         *global.Config
	clients     sync.Map
	clilen      int
	Receivers   sync.Map
	ethservice  *service.EthService
	addrservice *service.AddressService
}

func NewETHMonitor(cfg *global.Config, ethservice *service.EthService, addrservice *service.AddressService) *ETHMonitor {
	monitor := &ETHMonitor{
		cfg:         cfg,
		clients:     sync.Map{},
		Receivers:   sync.Map{},
		ethservice:  ethservice,
		addrservice: addrservice,
	}

	return monitor
}

func (m *ETHMonitor) newClient(chain global.ChainName) (client *ethclient.Client, err error) {
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

func (m *ETHMonitor) RangeListen() {
	type ListenRequest struct {
		Chain    string `json:"EVM chain"`
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

func (m *ETHMonitor) Listen(chain global.ChainName, oid string, receiver string, seconds int64) {
	m.Receivers.Store(receiver, true)
	defer m.Receivers.Delete(receiver)
	defer m.clearClients()

	item := &ClientItem{}
	found := false
	emu.Lock()
	m.clients.Range(func(key, val any) bool {
		cli := val.(*ClientItem)
		if cli.Chain == chain && cli.RunningQueryCount < MonitorClientSingleQueryLimit {
			item = cli
			found = true

			return false
		}

		return true
	})
	item.RunningQueryCount++
	emu.Unlock()

	if !found {
		for {
			if m.clilen < MonitorClientCount {
				client, err := m.newClient(chain)
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

	ichan := make(chan *entity.EthOrder, 1)
	go item.listen(ctx, ichan, sub, logs, receiver)

	for order := range ichan {
		order.MerchOrderId = oid
		order.Chain = string(chain)
		err = m.saveOrder(order)
		if err != nil {
			logx.Errorf("EVM chain 保存日志失败, txid:%v, err:%v", order.TxHash, err)
		}

		logx.Infof("🎉🎉🎉 chain 收到转账, [%v]:[%v], from:%v, to:%v", order.Currency, order.ReceivedAmount, order.FromHex, order.ToHex)

		err = global.NotifyEPay(m.cfg.EPay.NotifyUrl, order.MerchOrderId, order.FromHex, order.ToHex, order.Currency, order.ReceivedAmount)
		if err != nil {
			logx.Errorf("EVM chain 通知支付失败, moid:%v, txid:%v, err:%v", oid, order.TxHash, err)
			continue
		}

	}

	logx.Infof("EVM chain 事务结束, clen:%v, from:%v", m.clilen, receiver)
}

func (m *ETHMonitor) GenerateETHAddress() {
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

	m.ethservice.CreateAddresses(addrs)
}

func (m *ETHMonitor) clearClients() {
	emu.Lock()
	m.clients.Range(func(key, val any) bool {
		cli := val.(*ClientItem)
		if cli.RunningQueryCount <= 0 {
			cli.Client.Close()
			m.clients.Delete(key)
			m.clilen--
			logx.Infof("EVM chain 删除客户端, cname:%v, clen:%v", cli.Name, m.clilen)
		}

		return true
	})
	emu.Unlock()
}

func (m *ClientItem) listen(ctx context.Context, ichan chan *entity.EthOrder, sub ethereum.Subscription, logs chan types.Log, receiver string) {
	logx.Infof("EVM chain 实时状态开始, cname:%v, count:%v", m.Name, m.RunningQueryCount)
	defer func() {
		sub.Unsubscribe()
		close(logs)
		close(ichan)
		m.RunningQueryCount--

		logx.Infof("EVM chain 实时状态结束, unsub and close chans, cname:%v, count:%v", m.Name, m.RunningQueryCount)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("EVM chain 订阅超时, 已退出单笔订阅, to:%v", receiver)
			return
		case err := <-sub.Err():
			logx.Errorf("EVM chain 订阅错误, 已退出单笔订阅, to:%v, err:%v", receiver, err)
			return
		case log := <-logs:
			receipt, err1 := m.Client.TransactionReceipt(context.Background(), log.TxHash)
			if err1 != nil {
				logx.Errorf("EVM chain 获取交易回执失败: %s: %v", log.TxHash, err1)
				return
			}
			if receipt.Status != 1 {
				logx.Errorf("EVM chain 交易回执状态不为1, 可能已经挂起, txid:%s ", log.TxHash.String())
				return
			}

			to := common.HexToAddress(log.Topics[2].Hex()).Hex()
			if receiver == to {
				block, err := m.Client.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
				if err != nil {
					logx.Errorf("EVM chain log获取区块失败, txid:%v, err:%v", log.TxHash.String(), err)
					return
				}

				cid, err := m.Client.ChainID(context.Background())
				if err != nil {
					logx.Errorf("EVM chain 获取链ID失败, txid:%v, err:%v", log.TxHash.String(), err)
					return
				}

				from := common.HexToAddress(log.Topics[1].Hex())
				to := common.HexToAddress(log.Topics[2].Hex())

				sun := new(big.Int).SetBytes(log.Data)
				amount := global.Amount(sun.Int64(), global.AmountTypoBsc)

				var currency global.CurrenyTypo
				contract := strings.ToUpper(log.Address.Hex())
				for _, addr := range m.Cfg.ContractAddresses {
					if strings.ToUpper(addr.Address) == contract {
						currency = global.CurrenyTypo(addr.Currency)
						break
					}
				}

				order := &entity.EthOrder{
					Typo:           string(global.BscTransactionTypoIn),
					Status:         string(global.BscTransactionStatusSuccess),
					Currency:       string(currency),
					ChainId:        cid.Uint64(),
					TxHash:         log.TxHash.Hex(),
					TxIndex:        log.TxIndex,
					ReceivedAmount: amount,
					ReceivedSun:    sun.Int64(),
					FromHex:        from.Hex(),
					ToHex:          to.Hex(),
					Index:          log.Index,
					Contract:       log.Address.Hex(),
					BlockHash:      log.BlockHash.Hex(),
					BlockNumber:    log.BlockNumber,
					BlockTimestamp: block.Time(),
					Removed:        log.Removed,
					Description:    "",
				}

				ichan <- order
			}

			return
		}
	}
}

func (m *ETHMonitor) saveOrder(order *entity.EthOrder) (err error) {
	err = m.ethservice.SaveLog(order)
	if err != nil {
		logx.Errorf("EVM chain 记录转账交易失败: err:%v \n", err)
		return
	}

	return
}
