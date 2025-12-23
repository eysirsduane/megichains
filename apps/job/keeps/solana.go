package keeps

import (
	"context"
	"fmt"
	"log"
	"megichains/pkg/entity"
	"megichains/pkg/global"

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

func (m *SolanaClientItem) listen(ctx context.Context, ichan chan *entity.SolanaOrder, receiver string) {
	logx.Infof("SOLANA chain 实时状态开始, cname:%v, count:%v", m.Name, m.RunningQueryCount)
	defer func() {
		close(ichan)
		m.RunningQueryCount--

		logx.Infof("SOLANA chain 实时状态结束, unsub and close chans, cname:%v, count:%v", m.Name, m.RunningQueryCount)
	}()

	for {
		select {
		case <-ctx.Done():
			logx.Infof("SOLANA chain 订阅超时, 已退出单笔订阅, to:%v", receiver)
			return
		default:
			var lastAmount uint64 = 0

			for {
				var msg map[string]interface{}
				if err := m.Client.ReadJSON(&msg); err != nil {
					log.Println("read error:", err)
					continue
				}

				// 账户更新通知
				if msg["method"] == "accountNotification" {
					amount := parseUSDTAmount(msg)
					if lastAmount != 0 && amount != lastAmount {
						diff := int64(amount) - int64(lastAmount)
						fmt.Printf("💰 USDT Change: %+d (raw)\n", diff)
					}
					lastAmount = amount

					order := &entity.SolanaOrder{}

					ichan <- order
					
					break
				}
			}

		}

		return
	}
}

func parseUSDTAmount(msg map[string]interface{}) uint64 {
	params := msg["params"].(map[string]interface{})
	result := params["result"].(map[string]interface{})
	value := result["value"].(map[string]interface{})

	data := value["data"].(map[string]interface{})
	parsed := data["parsed"].(map[string]interface{})
	info := parsed["info"].(map[string]interface{})
	tokenAmount := info["tokenAmount"].(map[string]interface{})

	amountStr := tokenAmount["amount"].(string)

	var amount uint64
	fmt.Sscan(amountStr, &amount)
	return amount
}
