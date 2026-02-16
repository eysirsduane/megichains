package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/crypt"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service/clients"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

var emu sync.Mutex
var tmu sync.Mutex

const (
	MonitorClientCount             = 100
	MonitorClientSingleQueryLimit  = 20
	BitQueryClientSingleQueryLimit = 20
)

type ListenService struct {
	db            *gorm.DB
	cfg           *global.BackendesConfig
	clients       sync.Map
	clilen        int
	qclilen       int
	Receivers     sync.Map
	merchservice  *MerchService
	addrservice   *AddressService
	orderservice  *MerchOrderService
	chainservie   *ChainService
	evmservice    *EvmService
	tronservice   *TronService
	solanaservice *SolanaService
}

func NewListenService(cfg *global.BackendesConfig, db *gorm.DB, merchservice *MerchService, addrservice *AddressService, orderservice *MerchOrderService, chainservie *ChainService, evmservice *EvmService, tronservice *TronService, solanaservice *SolanaService) *ListenService {
	return &ListenService{
		cfg:           cfg,
		db:            db,
		clients:       sync.Map{},
		Receivers:     sync.Map{},
		merchservice:  merchservice,
		addrservice:   addrservice,
		orderservice:  orderservice,
		chainservie:   chainservie,
		evmservice:    evmservice,
		tronservice:   tronservice,
		solanaservice: solanaservice,
	}
}

func (s *ListenService) updateAddressLastUsed(addr string) {
	err := s.addrservice.SetLastUsed(addr)
	if err != nil {
		logx.Errorf("chain listen free address failed, addr:%v", addr)
		return
	}
}

func (s *ListenService) Listen(req *converter.ChainListenReq) {
	key := global.GetOrderAddressKey(string(req.Chain), req.Receiver, req.Currency)
	s.Receivers.Store(key, true)
	defer s.Receivers.Delete(key)
	defer s.clearClients()

	node, _ := snowflake.NewNode(1)
	sid := node.Generate()

	order := &entity.MerchantOrder{
		OrderNo:         sid.String(),
		MerchantAccount: req.MerchantAccount,
		MerchantOrderNo: req.MerchantOrderNo,
		Chain:           string(req.Chain),
		Typo:            string(global.OrderTypoIn),
		Status:          string(global.OrderStatusFailed),
		NotifyStatus:    string(global.NotifyStatusUnknown),
		Currency:        string(req.Currency),
		ToAddress:       req.Receiver,
		Description:     "",
	}
	err := s.orderservice.Save(order)
	if err != nil {
		err = biz.OrderSaveFailed
		logx.Errorf("chain listen order save failed, mono:%v, receiver:%v, err:%v", req.MerchantOrderNo, req.Receiver, err)
		return
	}

	s.updateAddressLastUsed(req.Receiver)

	req.Seconds += 120
	switch req.Chain {
	case global.ChainNameBsc, global.ChainNameEth:
		emu.Lock()
		c, err1 := s.getClientItem(req.Chain)
		if err1 != nil {
			logx.Errorf("EVM 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*clients.EvmClientItem)
		item.RunningQueryCount++
		emu.Unlock()

		s.listenEvm(order, req.Chain, global.CurrencyTypo(req.Currency), req.Receiver, req.Seconds, item)
	case global.ChainNameTron:
		tmu.Lock()
		c, err1 := s.getClientItem(req.Chain)
		if err1 != nil {
			logx.Errorf("TRON 获取客户端失败, 已退出事务.")
			return
		}
		item := c.(*clients.TronClientItem)
		tmu.Unlock()

		s.listenTron(order, req.Chain, global.CurrencyTypo(req.Currency), req.Receiver, req.Seconds, item)
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

		s.listenSolana(order, req.Chain, global.CurrencyTypo(req.Currency), req.Receiver, req.Seconds, item)
	default:
		logx.Errorf("=== 监听了未知的链 ===")
		return
	}

	log := &entity.MerchantOrderNotifyLog{
		MerchantOrderId: order.Id,
		NotifyUrl:       req.NotifyUrl,
	}
	err = s.orderservice.LogSave(log)
	if err != nil {
		logx.Errorf("chain listen order notify log save failed, mono:%v, txid:%v, err:%v", req.MerchantOrderNo, order.TransactionId, err)
		err = biz.OrderNotifyLogSaveFailed
		return
	}

	merch, err := s.merchservice.Get(req.MerchantAccount)
	if err != nil {
		order.Description = fmt.Sprintf("chain listen order notify not found merchant, account:%v", req.MerchantAccount)
		return
	}

	err = s.NotifyMerchant(log, merch, req.NotifyUrl, order.OrderNo, order.MerchantOrderNo, order.Status, order.TransactionId, order.FromAddress, order.ToAddress, order.Currency, order.ReceivedAmount, order.ReceivedSun)
	if err != nil {
		logx.Errorf("chain listen notify failed, chain:%v, ono:%v, mono:%v, txid:%v, err:%v", order.OrderNo, req.Chain, req.MerchantOrderNo, order.TransactionId, err)
		order.NotifyStatus = string(global.NotifyStatusFailed)
		order.Description = fmt.Sprintf("chain order notify failed. ono:%v, mono:%v, txid:%v, err:%v", order.OrderNo, req.MerchantOrderNo, order.TransactionId, err)

		log.Description = fmt.Sprintf("chain order notify failed. ono:%v, mono:%v, txid:%v, err:%v", order.OrderNo, req.MerchantOrderNo, order.TransactionId, err)
	} else {
		order.NotifyStatus = string(global.NotifyStatusSuccess)

	}

	err = s.orderservice.LogSave(log)
	if err != nil {
		logx.Errorf("chain listen order notify log update failed, ono:%v, mono:%v, txid:%v, err:%v", order.OrderNo, req.MerchantOrderNo, order.TransactionId, err)
		err = biz.OrderNotifyLogSaveFailed
	}

	err = s.orderservice.Save(order)
	if err != nil {
		logx.Errorf("chain listen order update failed, ono:%v, mono:%v, txid:%v, err:%v", order.OrderNo, req.MerchantOrderNo, order.TransactionId, err)
		err = biz.OrderSaveFailed
	}

	logx.Infof("监听事务结束, chain:%v, clen:%v, from:%v", req.Chain, s.clilen, req.Receiver)
}

func (s *ListenService) NotifyMerchant(log *entity.MerchantOrderNotifyLog, merch *entity.Merchant, url, ono, mono, status, txid, from, to, currency string, amount float64, sun int64) (err error) {
	req := global.OrderNotifyReq{
		OrderNo:         ono,
		MerchantOrderNo: mono,
		Status:          status,
		TxId:            txid,
		From:            from,
		To:              to,
		Currency:        currency,
		Amount:          amount,
		Sun:             sun,
	}

	rbytes := global.ObjToBytes(req)
	log.RequestBody = global.BytesToString(rbytes)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rbytes))
	if err != nil {
		return
	}

	sign := crypt.HmacSHA256(log.RequestBody, merch.SecretKey)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Merchant-Account", merch.MerchantAccount)
	request.Header.Add("Sign", sign)

	log.RequestHeader = global.ObjToJsonString(request.Header)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.ResponseHeader = global.ObjToJsonString(resp.Header)
	log.ResponseBody = global.BytesToString(bbytes)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("order notify response status is not ok, req:%+v, status:%v", req, resp.StatusCode)

		log.Description = err.Error()
		return
	}
	if log.RequestBody != "success" {
		err = fmt.Errorf("order notify response body is not success, req:%+v, body:%v", req, log.ResponseBody)

		log.Description = err.Error()
		return
	}

	return
}

func (s *ListenService) newEvmClient(chain global.ChainName) (client *ethclient.Client, err error) {
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

func (s *ListenService) newSolanaClient(chain global.ChainName) (wsc *ws.Client, rpcc *rpc.Client, err error) {
	wsc, err = ws.Connect(context.Background(), s.cfg.Solana.WssNetwork2)
	if err != nil {
		logx.Errorf("SOLANA new ws client failed, chain:%v, err:%v", chain, err)
		return
	}

	rpcc = rpc.New(s.cfg.Solana.GrpcNetwork)

	return
}

func (s *ListenService) newTronClient(chain global.ChainName) (conn *websocket.Conn, err error) {
	// c, err := client.NewClient(s.cfg.Tron.GrpcNetwork)
	// if err != nil {
	// 	logx.Errorf("Tron req.Chain connect failed, err:%v", err)
	// 	return
	// }

	logx.Infof("TRON new client, chain:%v", chain)

	return
}

func (s *ListenService) getClientItem(chain global.ChainName) (item any, err error) {
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

		tcli, ok := val.(*clients.TronClientItem)
		if ok {
			if tcli.Chain == chain && tcli.RunningQueryCount < MonitorClientSingleQueryLimit {
				item = tcli
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
					case global.ChainNameBsc, global.ChainNameEth:
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
						wsc, rpcc, err1 := s.newSolanaClient(chain)
						if err1 != nil {
							logx.Errorf("SOLANA chain Dial failed, try again err:%v", err1)
							time.Sleep(time.Second * 2)
							continue
						}

						c := &clients.SolanaClientItem{}
						name := uuid.NewString()
						c.Chain = chain
						c.Name = name
						c.WsClient = wsc
						c.RpcClient = rpcc
						c.Signatures = make(map[string]bool, 1)
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

func (s *ListenService) getContractAddress(currency global.CurrencyTypo, chain string) (caddr string, err error) {
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

func (s *ListenService) clearClients() {
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

		tcli, ok := val.(*clients.TronClientItem)
		if ok {
			if tcli.RunningQueryCount <= 0 {
				tcli.Client.Close()
				s.clients.Delete(key)
				s.clilen--
				logx.Infof("TRON chain delete client, cname:%v, clen:%v", tcli.Name, s.clilen)
			}
		}

		scli, ok := val.(*clients.SolanaClientItem)
		if ok {
			if scli.RunningQueryCount <= 0 {
				scli.WsClient.Close()
				scli.RpcClient.Close()
				s.clients.Delete(key)
				s.clilen--
				logx.Infof("SOLANA chain delete client, cname:%v, clen:%v", scli.Name, s.clilen)
			}

		}

		return true
	})

	emu.Unlock()
}

func (s *ListenService) listenEvm(order *entity.MerchantOrder, chain global.ChainName, currency global.CurrencyTypo, receiver string, seconds int64, item *clients.EvmClientItem) {
	contracts := make([]common.Address, 0, 1)
	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		logx.Errorf("EVM contract address get failed, chain:%v, currency:%v, err:%v", chain, currency, err)
		err = biz.ContractAddressNotFound
		return
	}
	contracts = append(contracts, common.HexToAddress(caddr))

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

	ichan := make(chan *entity.EvmLog, 1)
	go item.Listen(ctx, chain, ichan, currency, sub, logs, receiver)

	order.Status = string(global.OrderStatusTimeout)

	for log := range ichan {
		log.Chain = string(chain)
		err = s.evmservice.LogSave(log)
		if err != nil {
			logx.Errorf("evm service save log failed, err:%v", err)
			err = biz.EvmLogSaveFailed
			continue
		}

		logx.Infof("🎉🎉🎉 EVM 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, log.Currency, global.GetFloat64String(log.Amount), log.FromHex, log.ToHex)

		order.LogId = log.Id
		order.TransactionId = log.TxHash
		order.FromAddress = log.FromHex
		order.ReceivedAmount = log.Amount
		order.ReceivedSun = log.Sun
		order.Status = string(global.OrderStatusSuccess)

		s.chainservie.initChainClient(chain)
		s.chainservie.EvmFunds(receiver, chain)
	}

}

func (s *ListenService) listenTron(order *entity.MerchantOrder, chain global.ChainName, currency global.CurrencyTypo, receiver string, seconds int64, item *clients.TronClientItem) {
	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		logx.Errorf("Tron contract address get failed, chain:%v, currency:%v, err:%v", chain, currency, err)
		err = biz.ContractAddressNotFound
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	ichan := make(chan *entity.TronTransaction, 1)
	go item.Listen(ctx, chain, ichan, currency, s.cfg.Tron.HttpNetwork, caddr, receiver)

	order.Status = string(global.OrderStatusTimeout)

	for trans := range ichan {
		trans.Chain = string(chain)
		err = s.tronservice.TransSave(trans)
		if err != nil {
			logx.Errorf("tron service save trans failed, err:%v", err)
			err = biz.TronTransactionSaveFailed
			continue
		}

		logx.Infof("🎉🎉🎉 TRON 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, trans.Currency, global.GetFloat64String(trans.Amount), trans.FromBase58, trans.ToBase58)

		order.LogId = trans.Id
		order.TransactionId = trans.TransactionId
		order.Chain = string(chain)
		order.FromAddress = trans.FromBase58
		order.ReceivedAmount = trans.Amount
		order.ReceivedSun = trans.Sun
		order.Status = string(global.OrderStatusSuccess)

		s.chainservie.initChainClient(global.ChainNameTron)
		s.chainservie.TronFunds(receiver, global.ChainNameTron)
	}
}

func (s *ListenService) listenSolana(order *entity.MerchantOrder, chain global.ChainName, currency global.CurrencyTypo, receiver string, seconds int64, item *clients.SolanaClientItem) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	caddr, err := s.getContractAddress(currency, string(chain))
	if err != nil {
		logx.Errorf("SOLANA contract address get failed, chain:%v, currency:%v, err:%v", chain, currency, err)
		err = biz.ContractAddressNotFound
		return
	}
	umint := solana.MustPublicKeyFromBase58(caddr)

	raddr, err := solana.PublicKeyFromBase58(receiver)
	if err != nil {
		logx.Errorf("SOLANA parse base58 addrress failed, err:%v", err)
		return
	}

	ata, _, err := solana.FindAssociatedTokenAddress(raddr, umint)
	if err != nil {
		logx.Errorf("SOLANA find token address failed, raddr:%v, mint:%v, err:%v", raddr.String(), umint.String(), err)
		return
	}

	sub, err := item.WsClient.LogsSubscribeMentions(ata, rpc.CommitmentFinalized)
	if err != nil {
		logx.Errorf("SOLANA subscribe account transaction failed, receiver:%v, err:%v", receiver, err)
		return
	}

	ichan := make(chan *entity.SolanaTransaction, 1)
	go item.Listen(ctx, chain, ichan, currency, sub, receiver)

	order.Status = string(global.OrderStatusTimeout)

	for trans := range ichan {
		trans.Chain = string(chain)
		err = s.solanaservice.TransSave(trans)
		if err != nil {
			logx.Errorf("tron service save trans failed, err:%v", err)
			err = biz.SolanaTransactionSaveFailed
			continue
		}

		logx.Infof("🎉🎉🎉 SOLANA 收到转账, [%v]:[%v]:[%v], from:%v, to:%v", chain, trans.Currency, global.GetFloat64String(trans.Amount), trans.FromBase58, trans.ToBase58)

		order.LogId = trans.Id
		order.TransactionId = trans.TransactionId
		order.Chain = string(chain)
		order.FromAddress = trans.FromBase58
		order.ReceivedAmount = trans.Amount
		order.ReceivedSun = trans.Lamport
		order.Status = string(global.OrderStatusSuccess)

		s.chainservie.initChainClient(global.ChainNameSolana)
		s.chainservie.SolanaFunds(receiver, global.ChainNameSolana)

		return
	}
}

func (s *ListenService) ListenMany() {
	cnames := []string{"BSC", "ETH"}
	currencys := []string{"USDT", "USDC"}

	ctx := context.Background()

	for {
		addr, err := s.addrservice.FirstFree(ctx)
		if err != nil {
			break
		}

		if addr != nil {
			// if addr.Chain != string(global.ChainNameTron) {
			// 	continue
			// }

			icurr := rand.IntN(2)
			currency := currencys[icurr]
			if addr.Chain == string(global.ChainNameEvm) {
				random := rand.IntN(2)
				addr.Chain = cnames[random]
			}

			go s.Listen(&converter.ChainListenReq{
				MerchantOrderNo: uuid.NewString(),
				Chain:           global.ChainName(addr.Chain),
				Currency:        currency,
				Receiver:        addr.Address,
				Seconds:         600,
			})

			time.Sleep(time.Millisecond * 500)
		}
	}
}
