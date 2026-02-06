package erc20

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func PackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	// 从生成的合约代码中获取ABI（Erc20ABI 是abigen自动生成的全局变量）
	contractABI, err := abi.JSON(strings.NewReader(Erc20ABI))
	if err != nil {
		return nil, err
	}
	// 打包 transfer 方法的参数：地址 + 金额
	return contractABI.Pack("transfer", to, value)
}
