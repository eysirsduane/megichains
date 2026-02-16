package global

import "github.com/zeromicro/go-zero/rest"

type BackendesConfig struct {
	rest.RestConf
	Auth              *Auth
	DB                *DB
	Megi              *Megi
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

type Megi struct {
	NotifyUrl string
}

type Eth struct {
	ChainId     int64
	WssNetwork  string
	WssNetwork2 string
	ApiKey      string
}

type Bsc struct {
	ChainId     int64
	WssNetwork  string
	WssNetwork2 string
	ApiKey      string
}

type Solana struct {
	ChainId     int64
	WssNetwork  string
	WssNetwork2 string
	GrpcNetwork string
	ApiKey      string
}

type Tron struct {
	WssNetwork  string
	HttpNetwork string
	GrpcNetwork string
	ApiKey      string
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
