package global

import (
	"context"
	"fmt"
	"math/big"
	"megichains/pkg/biz"
	"megichains/pkg/erc20"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type AmountTypo int64

const (
	AmountTypoEth  AmountTypo = 1_000_000
	AmountTypoBsc  AmountTypo = 1_000_000_000_000_000_000
	AmountTypoSun  AmountTypo = 1_000_000
	AmountTypoTron AmountTypo = 1_000_000
)

func Amount(sun int64, typo AmountTypo) (amount float64) {
	a := decimal.NewFromInt(sun)
	b := decimal.NewFromFloat(float64(typo))

	return a.Div(b).InexactFloat64()
}

func Sun(amount float64, typo AmountTypo) (sun int64) {
	a := decimal.NewFromFloat(amount)
	b := decimal.NewFromFloat(float64(typo))
	return a.Mul(b).IntPart()
}

func TimeNowMilli() (milli uint64) {
	return uint64(time.Now().UnixMilli())
}

func TimeNowSeconds() (milli uint64) {
	return uint64(time.Now().Unix())
}

func TimeMaxMilli() (milli uint64) {
	return uint64(time.Now().AddDate(0, 0, 1).UnixMilli())
}

func TimeYesterdayMilli() (milli int64) {
	return time.Now().AddDate(0, 0, -1).UnixMilli()
}

func TimeYesterdaySeconds() (milli int64) {
	return time.Now().AddDate(0, 0, -1).Unix()
}

func TimeTomorrowSeconds() (milli int64) {
	return time.Now().AddDate(0, 0, 1).Unix()
}

func TimeTomorrowMilli() (milli uint64) {
	return uint64(time.Now().AddDate(0, 0, 1).UnixMilli())
}

func TimeMilliToSeconds(milli int64) (seconds int64) {
	return int64(milli / 1000)
}

func TimeInHourMilli() (milli uint64) {
	return uint64(time.Now().Add(time.Hour).UnixMilli())
}

func TimeLastHourMilli() (milli uint64) {
	return uint64(time.Now().Add(-time.Hour).UnixMilli())
}

func GetFloat64String(f float64) (s string) {
	return strconv.FormatFloat(f, 'f', 10, 64)
}

func GetOrderAddressKey(chain, receiver, currency string) (key string) {
	key = fmt.Sprintf("%v-%v-%v", chain, receiver, currency)
	return
}

func EthToWei(eth float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(eth), big.NewFloat(1e18))
	weiInt := new(big.Int)
	wei.Int(weiInt)
	return weiInt
}

func GweiToWei(gwei *big.Int) *big.Int {
	return new(big.Int).Mul(gwei, big.NewInt(1e9))
}

func UsdcToBase(usdc float64) *big.Int {
	f := new(big.Float).Mul(big.NewFloat(usdc), big.NewFloat(1e6))
	n := new(big.Int)
	f.Int(n)
	return n
}

func EstimateTransactionFee(ctx context.Context, cli *ethclient.Client, uaddr, fromAddress, toAddress common.Address, sun int64) (efee *big.Int, tcap *big.Int, fcap *big.Int, glimit uint64, err error) {
	// glimit = 80000
	txlimit, err := erc20.PackTransfer(toAddress, big.NewInt(sun))
	if err != nil {
		logx.Errorf("evm collect pack transfer for estimate gas limit failed, err:%v", err)
		err = biz.AddressFundCollectPackTransferFailed
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
		logx.Errorf("=== evm collect estimate gas limit failed ===, err:%v", err)
		glimit = 80000
	}
	glimit = uint64(float64(glimit) * 1.25)

	tcap, err = cli.SuggestGasTipCap(ctx)
	if err != nil {
		return
	}
	header, err := cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return
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

	return
}

func CheckFeeOverLimit(efee *big.Int, fmax float64) (bool, error) {
	wei := new(big.Float).Mul(
		new(big.Float).SetFloat64(fmax),
		new(big.Float).SetInt64(1e18),
	)

	mfee := new(big.Int)
	wei.Int(mfee)

	logx.Infof("check fee over limit, efee:%v, mfee:%v, fmax:%v", efee.String(), mfee.String(), fmax)
	return efee.Cmp(mfee) > 0, nil
}
