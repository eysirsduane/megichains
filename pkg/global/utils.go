package global

import (
	"time"

	"github.com/shopspring/decimal"
)

type AmountTypo int64

const (
	AmountTypoBsc AmountTypo = 1_000_000_000_000_000_000
	AmountTypoSun AmountTypo = 1_000_000
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
