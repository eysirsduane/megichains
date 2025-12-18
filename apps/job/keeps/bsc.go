package keeps

import (
	"context"
	"fmt"
	"math/big"
	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"megichains/pkg/service"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BSCMonitor struct {
	cfg        *global.Config
	bscservice *service.BscService
}

func NewBSCMonitor(cfg *global.Config, bscservice *service.BscService) *BSCMonitor {
	return &BSCMonitor{cfg: cfg, bscservice: bscservice}
}

const (
	BSC_WS   = "wss://bnb-testnet.g.alchemy.com/v2/e9NHfWClRByVJCAgbZsq7"
	Receiver = "0x0D0Aa17439878449c3E0f5211961738fEb35Fe37"
	Payer    = "0x9867551e5B82BA540A701387feD5e4F75Bc161d8"
)

func (m *BSCMonitor) newClient() (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(fmt.Sprintf("%v%v", m.cfg.Bsc.GrpcNetwork, m.cfg.Bsc.ApiKey))
	if err != nil {
		logx.Errorf("BSC Dial 失败 err:%v", err)
		return
	}

	return
}

func (m *BSCMonitor) Monitor() {
	for {
		client, err := m.newClient()
		if err != nil {
			logx.Errorf("BSC Dial 失败 err:%v", err)
			continue
		}

		usdtAddr := common.HexToAddress(m.cfg.Bsc.USDTAddress)
		usdcAddr := common.HexToAddress(m.cfg.Bsc.USDCAddress)
		receiver := common.HexToAddress(Receiver)
		// payer := common.HexToAddress(Payer)
		transferSig := crypto.Keccak256Hash(
			[]byte("Transfer(address,address,uint256)"),
		)

		query := ethereum.FilterQuery{
			Addresses: []common.Address{usdtAddr, usdcAddr},
			Topics: [][]common.Hash{
				{transferSig},
				nil,
				{common.HexToHash(receiver.Hex())},
			},
		}

		logs := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			logx.Errorf("BSC 订阅失败:%v", err)
			continue
		}

		fmt.Println("BSC 开始监听 BSC USDT USDC 转账...")

		for {
			select {
			case serr := <-sub.Err():
				logx.Errorf("BSC 订阅错误, 尝试重新连接, err:%v", serr)
				time.Sleep(time.Second * 5)
				client, err = m.newClient()
				if err != nil {
					logx.Errorf("BSC NewClient err:%v", err)
				}
			case log := <-logs:
				receipt, err1 := client.TransactionReceipt(context.Background(), log.TxHash)
				if err1 != nil {
					logx.Errorf("BSC 获取交易回执失败: %s: %v", log.TxHash, err1)
					continue
				}
				if receipt.Status != 1 {
					logx.Errorf("BSC 交易回执状态不为1, 可能已经挂起, txid:%s ", log.TxHash.String())
					continue
				}

				from := common.HexToAddress(log.Topics[1].Hex()).Hex()
				to := common.HexToAddress(log.Topics[2].Hex()).Hex()
				if Receiver == to {
					block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
					if err != nil {
						logx.Errorf("BSC log获取区块失败, txid:%v, err:%v", log.TxHash.String(), err)
						continue
					}

					currency, amount, err := m.handleLog(block.Time(), &log)
					if err != nil {
						logx.Errorf("BSC 处理log失败, txid:%v, err:%v", log.TxHash.String(), err)
						continue
					}

					logx.Infof("🎉🎉🎉 BSC 收到转账, [%v]:[%v], txid:%s, from:%v, to:%v", currency, amount, log.TxHash.String(), from, to)
				}
			}
		}
	}
}

func (m *BSCMonitor) handleLog(timestamp uint64, log *types.Log) (currency global.CurrenyTypo, amount float64, err error) {
	from := common.HexToAddress(log.Topics[1].Hex())
	to := common.HexToAddress(log.Topics[2].Hex())

	sun := new(big.Int).SetBytes(log.Data)
	amount = global.Amount(sun.Int64(), global.AmountTypoBsc)

	switch log.Address.Hex() {
	case m.cfg.Bsc.USDTAddress:
		currency = global.CurrencyTypoUsdt
	case m.cfg.Bsc.USDCAddress:
		currency = global.CurrencyTypoUsdc
	default:
		logx.Errorf("BSC 不支持的合约地址:%s", log.Address.Hex())
		err = biz.BscTransactionContractInvalid
		return
	}

	trans := &entity.BscTransaction{
		Typo:           string(global.BscTransactionTypoIn),
		Status:         string(global.BscTransactionStatusSuccess),
		Currency:       string(currency),
		ChainId:        m.cfg.Bsc.ChainId,
		TxHash:         log.TxHash.Hex(),
		TxIndex:        log.TxIndex,
		ReceivedAmount: amount,
		ReceivedSun:    sun.Int64(),
		FromHex:        from.Hex(),
		ToHex:          to.Hex(),
		Index:          log.Index,
		Address:        log.Address.Hex(),
		BlockHash:      log.BlockHash.Hex(),
		BlockNumber:    log.BlockNumber,
		BlockTimestamp: timestamp,
		Removed:        log.Removed,
		Description:    "",
	}

	err = m.bscservice.SaveTransaction(trans)
	if err != nil {
		logx.Errorf("BSC 记录转账日志失败: err:%v \n", err)
		return
	}

	return
}
