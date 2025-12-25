package keeps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type SolanaClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *websocket.Conn
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}

type LogsSubscribeReq struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type LogsNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				Signature string   `json:"signature"`
				Err       any      `json:"err"`
				Logs      []string `json:"logs"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

type TxResp struct {
	Result struct {
		Slot      uint64 `json:"slot"`
		BlockTime int64  `json:"blockTime"`
		Meta      struct {
			PreTokenBalances  []TokenBalance `json:"preTokenBalances"`
			PostTokenBalances []TokenBalance `json:"postTokenBalances"`
		} `json:"meta"`
	} `json:"result"`
}

type TokenBalance struct {
	Mint          string `json:"mint"`
	Owner         string `json:"owner"`
	UiTokenAmount struct {
		Amount   string `json:"amount"`
		Decimals uint8  `json:"decimals"`
	} `json:"uiTokenAmount"`
}

func (m *SolanaClientItem) listen(ctx context.Context, rpcurl, mint string, ichan chan *entity.SolanaOrder, receiver string) {
	logx.Infof("SOLANA chain 实时状态开始, cname:%v, count:%v, receiver:%v", m.Name, m.RunningQueryCount, receiver)
	defer func() {
		close(ichan)
		m.RunningQueryCount--
		logx.Infof("SOLANA chain 实时状态结束, unsub and close chans, cname:%v, count:%v, receiver:%v", m.Name, m.RunningQueryCount, receiver)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("SOLANA chain 订阅超时, 已退出单笔订阅, to:%v", receiver)
			return
		default:
			for {
				var log LogsNotification
				if err := m.Client.ReadJSON(&log); err != nil {
					logx.Errorf("read err:%v", err)
					return
				}

				if log.Method != "logsNotification" {
					continue
				}

				if log.Params.Result.Value.Err != nil {
					continue
				}

				sig := log.Params.Result.Value.Signature
				tx := m.getTransaction(rpcurl, sig)
				if tx == nil || tx.Result.Meta.PostTokenBalances == nil {
					return
				}

				pre := make(map[string]uint64)

				// 先记录 pre balances
				for _, b := range tx.Result.Meta.PreTokenBalances {
					if b.Mint == mint && b.Owner == receiver {
						amt, _ := strconv.ParseUint(b.UiTokenAmount.Amount, 10, 64)
						pre[b.Owner] = amt
					}
				}

				// 再对比 post balances
				for _, b := range tx.Result.Meta.PostTokenBalances {
					if b.Mint != mint || b.Owner != receiver {
						continue
					}

					postAmt, _ := strconv.ParseUint(b.UiTokenAmount.Amount, 10, 64)
					preAmt := pre[b.Owner]

					if postAmt > preAmt {
						diff := postAmt - preAmt

						fmt.Printf(`
								====== INCOMING TRANSFER ======
								Chain        : solana-devnet
								TxID         : %s
								Mint         : %s
								To           : %s
								Amount       : %d
								Slot         : %d
								BlockTime    : %d
								================================
								`,
							sig,
							mint,
							receiver,
							diff,
							tx.Result.Slot,
							tx.Result.BlockTime,
						)

						break
					}
				}
			}
		}
	}
}

func (m *SolanaClientItem) getTransaction(rpcurl string, signature string) *TxResp {
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTransaction",
		"params": []any{
			signature,
			map[string]string{
				"encoding": "jsonParsed",
			},
		},
	}

	b, _ := json.Marshal(payload)
	resp, err := http.Post(rpcurl, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var r TxResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil
	}

	return &r
}

func isUSDTTransfer(logs []string) bool {
	hasTokenProgram := false
	hasTransfer := false
	hasUSDT := false

	for _, l := range logs {
		if strings.Contains(l, "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA") {
			hasTokenProgram = true
		}
		if strings.Contains(l, "Instruction: Transfer") {
			hasTransfer = true
		}
		if strings.Contains(l, "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB") {
			hasUSDT = true
		}
	}

	return hasTokenProgram && hasTransfer && hasUSDT
}
