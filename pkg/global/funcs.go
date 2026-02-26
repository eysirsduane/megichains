package global

import (
	"encoding/json"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Offset(page, limit int) (offset int) {
	return (page - 1) * limit
}

func GetTrx2UsdtRateFromHtx(url string) (rate float64, err error) {

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
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username, secret string, expire int64, issuer string) (token string, err error) {
	claims := Claims{
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

func GenerateRefreshToken(username, secret string, expire int64, issuer string) (token string, err error) {
	claims := Claims{
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

func ObjToBytes(obj any) (bytes []byte) {
	bytes, _ = json.Marshal(obj)

	return
}

func ObjToJsonString(obj any) (str string) {
	bytes := ObjToBytes(obj)
	str = string(bytes)

	return
}

func BytesToString(bytes []byte) (str string) {
	str = string(bytes)
	return
}

func BytesToObj(bytes []byte, obj any) (err error) {
	err = json.Unmarshal(bytes, obj)
	return
}
