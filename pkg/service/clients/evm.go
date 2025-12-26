package clients

import (
	"context"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EvmClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *ethclient.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

func (m *EvmClientItem) Listen(ctx context.Context, chain global.ChainName, currency string, ichan chan *entity.EvmOrder, sub ethereum.Subscription, logs chan types.Log, receiver string) {
	logx.Infof("EVM chain 实时状态开始, cname:%v, count:%v, chain:%v, currency:%v, receiver:%v", m.Name, m.RunningQueryCount, chain, currency, receiver)
	defer func() {
		sub.Unsubscribe()
		close(logs)
		close(ichan)
		m.RunningQueryCount--

		logx.Infof("EVM chain 实时状态结束, unsub and close chans, cname:%v, count:%v, chain:%v, currency:%v, receiver:%v", m.Name, m.RunningQueryCount, chain, currency, receiver)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("EVM chain 订阅超时, 已退出单笔订阅, receiver:%v", receiver)
			return
		case err := <-sub.Err():
			logx.Errorf("EVM chain 订阅错误, 已退出单笔订阅, receiver:%v, err:%v", receiver, err)
			return
		case log := <-logs:
			if log.Removed {
				logx.Errorf("EVM log is removed, currency:%v, receiver:%v, txid:%v", currency, receiver, log.TxHash.String())
				return
			}

			round, threshold := 0, 0
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
					logx.Errorf("EVM 经过11轮比较区块高度后还未稳定, 认为本次交易失败... chain:%v, currency:%v, receiver:%v", chain, currency, receiver)
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
