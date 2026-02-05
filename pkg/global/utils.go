package global

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
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
