// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package global

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth Auth
	DB   DB
	Tron Tron
	Bsc  Bsc
	Bot  Bot
}

type Auth struct {
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
	Issuer        string
}

type DB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int16
	Charset  string
	Timezone string
}

type Bsc struct {
	ChainId     uint16
	GrpcNetwork string
	ApiKey      string
	USDTAddress string
	USDCAddress string
}

type Tron struct {
	GrpcNetwork string
	HttpNetwork string

	Trx2UsdtRateApi string

	OwnerPrivateKey  string
	MonitorAddress   string
	TRC20USDTAddress string
}

type Bot struct {
	Token   string
	Service string
}
