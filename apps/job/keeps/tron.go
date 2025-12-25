package keeps

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

func (m *TronClientItem) listen(ctx context.Context, ichan chan *entity.TronOrder, currency global.CurrencyTypo, httpurl, caddr, receiver string) {
	logx.Infof("TRON chain 实时状态开始, cname:%v, count:%v, receiver:%v", m.Name, m.RunningQueryCount, receiver)
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

								amount := global.Amount(sun, global.AmountTypoTron)

								order := &entity.TronOrder{}
								order.TransactionId = tx.TransactionId
								order.BlockTimestamp = tx.BlockTimestamp
								order.Contract = tx.TokenInfo.Address
								order.Currency = tx.TokenInfo.Symbol
								order.FromBase58 = tx.From
								order.ToBase58 = tx.To
								order.ReceivedAmount = amount
								order.ReceivedSun = sun
								order.Typo = string(global.BscTransactionTypoIn)
								order.Status = string(global.BscTransactionStatusSuccess)
								order.Description = ""

								ichan <- order

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

// func startConfirmationWorker() {
// 	ticker := time.NewTicker(5 * time.Second)

// 	for range ticker.C {
// 		latest, err := getLatestBlock()
// 		if err != nil {
// 			continue
// 		}

// 		mu.Lock()
// 		for _, tx := range pendingTxs {
// 			if tx.Confirmed {
// 				continue
// 			}

// 			if latest >= tx.BlockHeight+CONFIRMATIONS {
// 				tx.Confirmed = true

// 				log.Printf(`
// ✅ CONFIRMED USDT INCOMING
// TxID        : %s
// From        : %s
// To          : %s
// Amount      : %s
// BlockHeight : %d
// `,
// 					tx.TxID,
// 					tx.From,
// 					tx.To,
// 					tx.Amount,
// 					tx.BlockHeight,
// 				)

// 				// 👉 在这里做真正的入账逻辑
// 				// creditUser(tx)
// 			}
// 		}
// 		mu.Unlock()
// 	}
// }

// func getLatestBlock() (uint64, error) {
// 	resp, err := http.Get(TRON_RPC + "/wallet/getnowblock")
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()

// 	var r struct {
// 		BlockHeader struct {
// 			RawData struct {
// 				Number uint64 `json:"number"`
// 			} `json:"raw_data"`
// 		} `json:"block_header"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
// 		return 0, err
// 	}

// 	return r.BlockHeader.RawData.Number, nil
// }

// var mu sync.Mutex
// var pendingTxs = make(map[string]*PendingTx)

// type PendingTx struct {
// 	TxID        string
// 	From        string
// 	To          string
// 	Amount      string
// 	BlockHeight uint64
// 	Confirmed   bool
// }

// type BitqueryMsg struct {
// 	Type    string `json:"type"`
// 	Payload struct {
// 		Data struct {
// 			Tron struct {
// 				Transfers []struct {
// 					Transaction struct {
// 						Hash string `json:"hash"`
// 					}
// 					Amount string `json:"amount"`
// 					Sender struct {
// 						Address string `json:"address"`
// 					}
// 					Receiver struct {
// 						Address string `json:"address"`
// 					}
// 					Block struct {
// 						Height uint64 `json:"height"`
// 					}
// 				} `json:"transfers"`
// 			} `json:"tron"`
// 		} `json:"data"`
// 	} `json:"payload"`
// }

// func handleBitqueryMsg(raw []byte) {
// 	var msg BitqueryMsg
// 	if err := json.Unmarshal(raw, &msg); err != nil {
// 		return
// 	}

// 	if msg.Type != "data" {
// 		return
// 	}

// 	for _, t := range msg.Payload.Data.Tron.Transfers {
// 		mu.Lock()
// 		if _, exists := pendingTxs[t.Transaction.Hash]; exists {
// 			mu.Unlock()
// 			continue
// 		}

// 		pendingTxs[t.Transaction.Hash] = &PendingTx{
// 			TxID:        t.Transaction.Hash,
// 			From:        t.Sender.Address,
// 			To:          t.Receiver.Address,
// 			Amount:      t.Amount,
// 			BlockHeight: t.Block.Height,
// 		}
// 		mu.Unlock()

// 		log.Printf("🕒 PENDING tx=%s amount=%s block=%d", t.Transaction.Hash, t.Amount, t.Block.Height)
// 	}
// }
