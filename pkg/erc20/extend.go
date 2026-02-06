package erc20

import (
	"context"
	"math/big"
	"megichains/pkg/global"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func PackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	contractABI, err := abi.JSON(strings.NewReader(Erc20ABI))
	if err != nil {
		return nil, err
	}

	return contractABI.Pack("transfer", to, value)
}

func EstimateTransactionFee(ctx context.Context, cli *ethclient.Client, chain global.ChainName, uaddr, fromAddress, toAddress common.Address, sun int64) (efee *big.Int, tcap *big.Int, fcap *big.Int, glimit uint64, err error) {
	txlimit, err := PackTransfer(toAddress, big.NewInt(sun))
	if err != nil {
		return
	}
	msg := ethereum.CallMsg{
		From:  fromAddress,   // 付款账户
		To:    &uaddr,        // 目标合约地址
		Data:  txlimit,       // 合约调用数据
		Value: big.NewInt(0), // ERC20转账不转账ETH，值为0
	}
	glimit, err = cli.EstimateGas(ctx, msg)
	if err != nil {
		glimit = 80000
	}
	glimit = uint64(float64(glimit) * 1.25)

	switch chain {
	case global.ChainNameBsc:
		fcap, err := cli.SuggestGasPrice(ctx)
		if err != nil {
			break
		}

		efee = new(big.Int).Mul(fcap, big.NewInt(int64(glimit)))
	case global.ChainNameEth:
		tcap, err = cli.SuggestGasTipCap(ctx)
		if err != nil {
			break
		}
		header, err := cli.HeaderByNumber(ctx, nil)
		if err != nil {
			break
		}

		fcap = new(big.Int).Add(header.BaseFee, tcap)

		mtip := big.NewInt(40_000_000)
		mfcap := big.NewInt(400_000_000)
		if tcap.Cmp(mtip) < 0 {
			tcap = mtip
		}
		if fcap.Cmp(mfcap) < 0 {
			fcap = mfcap
		}

		efee = new(big.Int).Mul(fcap, big.NewInt(int64(glimit)))
	}

	return
}

func CheckFeeOverLimit(efee *big.Int, fmax float64) (bool, error) {
	wei := new(big.Float).Mul(
		new(big.Float).SetFloat64(fmax),
		new(big.Float).SetInt64(1e18),
	)

	mfee := new(big.Int)
	wei.Int(mfee)

	return efee.Cmp(mfee) > 0, nil
}
