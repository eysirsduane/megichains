package clients

import (
	"context"
	"math/big"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"strings"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/zeromicro/go-zero/core/logx"
)

type SolanaClientItem struct {
	Name              string
	Chain             global.ChainName
	WsClient          *ws.Client
	RpcClient         *rpc.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
	Signatures        map[string]bool
}

func (m *SolanaClientItem) Listen(ctx context.Context, chain global.ChainName, ichan chan *entity.SolanaTransaction, currency global.CurrencyTypo, sub *ws.LogSubscription, receiver string) {
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
		case msg := <-sub.Response():
			ctx := context.Background()
			if msg.Value.Err != nil {
				logx.Errorf("SOLANA sub response value has error, receiver:%v, err:%v", receiver, msg.Value.Err)
				return
			}

			sig := msg.Value.Signature
			if m.Signatures[sig.String()] == true {
				return
			}

			if msg.Value.Signature.IsZero() {
				continue
			}

			if !containsTransfer(msg.Value.Logs) {
				continue
			}

			tx, err := m.RpcClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{Commitment: rpc.CommitmentConfirmed})
			if err != nil {
				logx.Errorf("SOLANA get transaction failed, receiver:%v, err:%v", receiver, err)
				return
			}

			if tx.Meta == nil {
				logx.Errorf("SOLANA tx meta is nil, receiver:%v", receiver)
				return
			}

			from := ""
			ttx, _ := tx.Transaction.GetTransaction()
			accs, _ := ttx.AccountMetaList()
			for _, addr := range accs {
				if addr.IsSigner {
					from = addr.PublicKey.String()
					logx.Infof("SOLANA trans account meta list get from address, currency:%v, from:%v, receiver:%v, txid:%v", currency, from, receiver, sig.String())

					break
				}
			}

			mint := ""
			amount := float64(0)
			sun := int64(0)
			for _, post := range tx.Meta.PostTokenBalances {
				mint = post.Mint.String()
				if post.Owner.String() == receiver {
					for _, pre := range tx.Meta.PreTokenBalances {
						if pre.AccountIndex == post.AccountIndex {
							preAmt, _ := new(big.Int).SetString(pre.UiTokenAmount.Amount, 10)
							postAmt, _ := new(big.Int).SetString(post.UiTokenAmount.Amount, 10)
							diff := new(big.Int).Sub(postAmt, preAmt)
							sun = diff.Int64()

							amount = *post.UiTokenAmount.UiAmount - *pre.UiTokenAmount.UiAmount
						}
					}
				}
			}

			trans := &entity.SolanaTransaction{
				Currency:      string(currency),
				Chain:         string(chain),
				TransactionId: sig.String(),
				Amount:        amount,
				Lamport:       sun,
				FromBase58:    from,
				ToBase58:      receiver,
				Mint:          mint,
				BlockTime:     uint64(*tx.BlockTime),
			}

			m.Signatures[sig.String()] = true

			ichan <- trans
		}
	}
}

func containsTransfer(logs []string) bool {
	for _, l := range logs {
		if strings.Contains(l, "Instruction: Transfer") {
			return true
		}
	}
	return false
}
