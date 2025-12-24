// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package global

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth              *Auth
	DB                *DB
	EPay              *EPay
	Tron              *Tron
	Eth               *Eth
	Bsc               *Bsc
	Solana            *Solana
	Bot               *Bot
	ContractAddresses []*ContractAddress
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

type EPay struct {
	NotifyUrl string
}

type Eth struct {
	ChainId     uint16
	GrpcNetwork string
	ApiKey      string
}

type Bsc struct {
	ChainId     uint16
	GrpcNetwork string
	ApiKey      string
}

type Solana struct {
	ChainId     uint16
	GrpcNetwork string
	ApiKey      string
	UsdtMint    string
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

type ContractAddress struct {
	Chain    string
	Address  string
	Currency string
}

type ListenReq struct {
	MerchOrderId string `json:"merch_order_id"`
	Chain        string `json:"chain"`
	Currency     string `json:"currency"`
	Receiver     string `json:"receiver"`
	Seconds      int64  `json:"seconds"`
}

type EPayNotifyReq struct {
	MerchOrderId string  `json:"merch_order_id"`
	TxId         string  `json:"txid"`
	FromHex      string  `json:"from_hex"`
	ToHex        string  `json:"to_hex"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
}

type EPayOrderAddress struct {
	Address  string
	Currency string
}
