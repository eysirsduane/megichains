package keeps

import (
	"megichains/pkg/global"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	MonitorClientCount            = 100
	MonitorClientSingleQueryLimit = 20
)

type ClientItem struct {
	Cfg               *global.Config
	Name              string
	Chain             global.ChainName
	Client            *ethclient.Client
	Status            int // 0: 空闲, 1: 使用中
	RunningQueryCount int
}
