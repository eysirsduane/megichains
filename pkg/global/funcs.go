package global

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"megichains/pkg/biz"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

func Offset(page, limit int) (offset int) {
	return (page - 1) * limit
}

func GetTrx2UsdtRateFromHtx(url string) (rate float64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		logx.Infof("funcs get trx2usdt rate http request failed, err:%v", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	wrapper := &struct {
		Status string `json:"status"`
		Tick   struct {
			Open float64 `json:"open"`
		} `json:"tick"`
	}{}

	if err = json.Unmarshal(body, wrapper); err != nil {
		logx.Infof("funcs get trx2usdt rate unmarshal json failed, err:%v", err)
		return
	}

	rate = wrapper.Tick.Open
	if rate >= 1 || rate <= 0 {
		logx.Errorf("funcs get trx2usdt rate value is incorrect, rate:%v", rate)
		err = biz.ExchangeTrx2UsdtRateIncorrect
		return
	}

	return
}

func GetTrx2UsdtExchangeRate(rate float64) (erate float64) {
	return math.Ceil(rate*100) / 100
}

func GetUsdt2TrxAmount(rate float64, amount float64, discount float64) (eamount, erate float64) {
	erate = GetTrx2UsdtExchangeRate(rate)
	eamount = math.Floor(((amount/erate)*discount)*100) / 100

	return
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(uid int64, username string, secret string, expire int64, issuer string) (token string, err error) {
	claims := Claims{
		uid,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   "megichains",
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tk.SignedString([]byte(secret))

	return
}

func GenerateRefreshToken(uid int64, username string, secret string, expire int64, issuer string) (token string, err error) {
	claims := Claims{
		UserID:   uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   "megichains",
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tk.SignedString([]byte(secret))

	return
}

func ParseToken(token string, secret string) (*Claims, error) {
	tk, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return tk.Claims.(*Claims), nil
}

func HashPassword(pwd string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	hash = string(bytes)
	return
}

func CheckPassword(hash, pwd string) (ok bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func NotifyEPay(url, merchOrderId, fromHex, toHex, currency string, receivedAmount float64) (err error) {
	req := EPayRequest{
		MerchOrderId: merchOrderId,
		FromHex:      fromHex,
		ToHex:        toHex,
		Amount:       receivedAmount,
		Currency:     currency,
	}

	byts, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(byts))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notify epay service failed, status code:%v", resp.StatusCode)
	}

	return
}
