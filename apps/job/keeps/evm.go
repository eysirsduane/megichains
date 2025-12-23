package keeps

import (
	"context"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var emu sync.Mutex

type EvmClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *ethclient.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

func (m *EvmClientItem) listen(ctx context.Context, chain global.ChainName, ichan chan *entity.EvmOrder, sub ethereum.Subscription, logs chan types.Log, receiver string) {
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
			// for i := 0; i < 5; i++ {
			// 	receipt, err1 := m.Client.TransactionReceipt(context.Background(), log.TxHash)
			// 	if err1 != nil {
			// 		logx.Errorf("EVM chain 获取交易回执失败: %s: %v", log.TxHash, err1)
			// 		continue
			// 	}
			// 	if receipt.Status != 1 {
			// 		logx.Errorf("EVM chain 交易回执状态不为1, 可能已经挂起, txid:%s ", log.TxHash.String())
			// 		continue
			// 	}
			// }
			round := 0
			threshold := 0
			switch chain {
			case global.ChainNameEth:
				threshold = 1
			case global.ChainNameBsc:
				threshold = 7
			default:
				logx.Errorf("EVM 计算区块阈值失败, 不支持这个链, chain:%v", chain)
				return
			}
			for {
				curblock, err1 := m.Client.BlockByNumber(context.Background(), nil)
				if err1 != nil {
					logx.Errorf("EVM chain 获取最新区块发生错误, err:%v", err1)
					time.Sleep(time.Second * 5)
					continue
				}

				logx.Infof("EVM 获取区块高度, cur:%v, log:%v", curblock.NumberU64(), log.BlockNumber)

				if curblock.NumberU64()-log.BlockNumber >= uint64(threshold) {
					break
				}

				if round > 10 {
					logx.Errorf("EVM 经过11轮比较区块高度后还未稳定, 所以认为本次交易失败...!")
					return
				}

				time.Sleep(time.Second * 5)
				round++
			}

			to := common.HexToAddress(log.Topics[2].Hex()).Hex()
			if strings.EqualFold(receiver, to) {
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
				amount := float64(0)
				switch chain {
				case global.ChainNameEth:
					amount = global.Amount(sun.Int64(), global.AmountTypoEth)
				case global.ChainNameBsc:
					amount = global.Amount(sun.Int64(), global.AmountTypoBsc)
				default:
					logx.Errorf("EVM chain 未知链SUN/AMOUNT比值, chain:%v", chain)
					return
				}

				var currency global.CurrencyTypo
				contract := strings.ToUpper(log.Address.Hex())
				for _, addr := range m.Cfg.ContractAddresses {
					if strings.EqualFold(addr.Address, contract) {
						currency = global.CurrencyTypo(addr.Currency)
						break
					}
				}

				order := &entity.EvmOrder{
					Typo:           string(global.BscTransactionTypoIn),
					Status:         string(global.BscTransactionStatusSuccess),
					Currency:       string(currency),
					ChainId:        cid.Uint64(),
					TxHash:         log.TxHash.Hex(),
					Index:          log.Index,
					TxIndex:        log.TxIndex,
					ReceivedAmount: amount,
					ReceivedSun:    sun.Int64(),
					FromHex:        from.Hex(),
					ToHex:          to.Hex(),
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
