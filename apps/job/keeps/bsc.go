package keeps

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type BSCMonitor struct {
	db *gorm.DB
}

func NewBSCMonitor(db *gorm.DB) *BSCMonitor {
	return &BSCMonitor{
		db: db,
	}
}

const (
	BSC_WS   = "wss://bnb-testnet.g.alchemy.com/v2/e9NHfWClRByVJCAgbZsq7"
	USDT     = "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd"
	USDC     = "0x64544969ed7ebf5f083679233325356ebe738930"
	Receiver = "0x0D0Aa17439878449c3E0f5211961738fEb35Fe37"
	Payer    = "0x9867551e5B82BA540A701387feD5e4F75Bc161d8"
)

func NewClient() (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(BSC_WS)
	if err != nil {
		fmt.Println("BSC dial bsc err:", err)
		return
	}

	return
}

func (m *BSCMonitor) Monitor() {
	client, err := NewClient()
	if err != nil {
		fmt.Println("BSC NewClient err:", err)
		return
	}

	usdtAddr := common.HexToAddress(USDT)
	usdcAddr := common.HexToAddress(USDC)
	receiver := common.HexToAddress(Receiver)
	payer := common.HexToAddress(Payer)
	transferSig := crypto.Keccak256Hash(
		[]byte("Transfer(address,address,uint256)"),
	)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{usdtAddr, usdcAddr},
		Topics: [][]common.Hash{
			{transferSig},
			nil,
			{common.HexToHash(receiver.Hex()), common.HexToHash(payer.Hex())},
		},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		fmt.Println("BSC 订阅失败:", err)
		return
	}

	fmt.Println("BSC 开始监听 BSC USDT 转账...")

	for {
		select {
		case serr := <-sub.Err():
			fmt.Println("BSC 订阅错误:", serr)
			time.Sleep(time.Second * 5)
			client, err = NewClient()
			if err != nil {
				fmt.Println("BSC NewClient err:", err)
			}
		case log := <-logs:
			txInfo, err1 := client.TransactionReceipt(context.Background(), log.TxHash)
			if err1 != nil {
				fmt.Printf("BSC get transaction receipt for %s: %v", log.TxHash, err1)
				continue
			}
			if txInfo.Status != 1 {
				fmt.Printf("BSC transaction isPending txid %s ", log.TxHash.String())
				continue
			}

			to := common.HexToAddress(log.Topics[2].Hex())
			if Receiver == to.Hex() {
				handleLog(log)
			}

		}
	}
}

func handleLog(log types.Log) {
	txTime := time.Unix(int64(log.BlockTimestamp), 0)

	from := common.HexToAddress(log.Topics[1].Hex())
	to := common.HexToAddress(log.Topics[2].Hex())

	// USDT decimals = 18 (BSC)
	amount := new(big.Int).SetBytes(log.Data)
	contract := log.Address.Hex()

	fmt.Printf("BSC 接收转账: 合约=%s, 来自=%s, 去往=%s, 数量=%s, 区块=%d, 时间=%s, 交易=%s, 块:=%s\n",
		contract,
		from.Hex(),
		to.Hex(),
		amount.String(),
		log.BlockNumber,
		txTime.Format("2006-01-02 15:04:05"),
		log.TxHash.Hex(),
		log.BlockHash.String(),
	)

	fmt.Printf("log details:%+v \n", log)
}
