package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type TronClientItem struct {
	Name              string
	Chain             global.ChainName
	Client            *websocket.Conn
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

func (m *TronClientItem) Listen(ctx context.Context, chain global.ChainName, ichan chan *entity.TronTransaction, currency global.CurrencyTypo, httpurl, caddr, receiver string) {
	logx.Infof("TRON chain 实时状态开始, chain:%v, currency:%v, receiver:%v, cname:%v, count:%v, ", chain, currency, receiver, m.Name, m.RunningQueryCount)
	defer func() {
		close(ichan)
		m.RunningQueryCount--
		logx.Infof("TRON chain 实时状态结束, close chans, cname:%v, count:%v, receiver:%v", m.Name, m.RunningQueryCount, receiver)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("TRON chain 订阅超时, 已退出单笔订阅, currency:%v, to:%v", currency, receiver)
			return
		default:
			min := time.Now().UnixMilli()

			for {
				url := fmt.Sprintf("%v/v1/accounts/%s/transactions/trc20?limit=200&contract_address=%v&only_confirmed=true&min_timestamp=%v", httpurl, receiver, caddr, min)
				resp, err := http.Get(url)
				if err != nil {
					logx.Infof("Tron %v transaction 监听请求失败, err:%v", currency, err)
					return
				}

				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				wrapper := &struct{ Data []*Trc20Transaction }{}
				if err := json.Unmarshal(body, wrapper); err != nil {
					logx.Infof("Tron %v transaction 监听序列化json失败, err:%v", currency, err)
					return
				}

				if len(wrapper.Data) > 0 {
					for _, tx := range wrapper.Data {
						logx.Infof("Tron 获得监听消息, [%v], receiver:%v", currency, receiver)
						if tx != nil {
							if tx.Type != global.TronTransactionTypoTransfer {
								continue
							}
							if tx.From == "" || tx.To == "" {
								continue
							}

							if tx.To == receiver {
								sun, err := strconv.ParseInt(tx.Value, 10, 64)
								if err != nil {
									logx.Errorf("Tron transaction 转换金额是啊币, [%v]:[%v], receiver:%v, err:%v", currency, tx.Value, receiver, err)
									continue
								}

								amount := global.Amount(sun, global.AmountTypo6e)

								trans := &entity.TronTransaction{}
								trans.Currency = tx.TokenInfo.Symbol
								trans.TransactionId = tx.TransactionId
								trans.Amount = amount
								trans.Sun = sun
								trans.FromBase58 = tx.From
								trans.ToBase58 = tx.To
								trans.Contract = tx.TokenInfo.Address
								trans.BlockTimestamp = tx.BlockTimestamp

								min = int64(tx.BlockTimestamp)

								ichan <- trans

								return
							}

							min = int64(tx.BlockTimestamp)
						}
					}
				}

				time.Sleep(time.Second * 5)
			}
		}
	}
}

type Trc20Transaction struct {
	TransactionId string `json:"transaction_id"`
	TokenInfo     struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int32  `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	From           string `json:"from"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Value          string `json:"value"`
}
