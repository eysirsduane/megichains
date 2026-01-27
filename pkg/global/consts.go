package global

const (
	StringSuccess = "SUCCESS"
)

type AddressTypo string

const (
	AddressTypoIn      AddressTypo = "IN"
	AddressTypoOut     AddressTypo = "OUT"
	AddressTypoCollect AddressTypo = "COLLECT"
)

type AddressStatus string

const (
	AddressStatusBaned  AddressStatus = "禁用"
	AddressStatusInUse  AddressStatus = "占用"
	AddressStatusInFree AddressStatus = "空闲"
)

type OrderTypo string

const (
	OrderTypoIn  OrderTypo = "输入"
	OrderTypoOut OrderTypo = "输出"
)

const (
	TronTransactionTypoTransfer                 string = "Transfer"
	TronTransactionTypoTransferContract         string = "TransferContract"
	TronTransactionTypoDelegateResourceContract string = "DelegateResourceContract"
)

type ChainName string

const (
	ChainNameTron   ChainName = "TRON"
	ChainNameEth    ChainName = "ETH"
	ChainNameBsc    ChainName = "BSC"
	ChainNameEvm    ChainName = "EVM"
	ChainNameSolana ChainName = "SOLANA"
)

type CurrencyTypo string

const (
	CurrencyTypoTrx    CurrencyTypo = "TRX"
	CurrencyTypoUsdt   CurrencyTypo = "USDT"
	CurrencyTypoUsdc   CurrencyTypo = "USDC"
	CurrencyTypoEnergy CurrencyTypo = "ENERGY"
)

type ExchangeRateTypo string

const (
	ExchangeRateTrx2Energy ExchangeRateTypo = "TRX2ENERGY"
	ExchangeRateUsdt2Trx   ExchangeRateTypo = "USDT2TRX"
)

type DelegateResourceTypo uint32

const (
	DelegateResourceTypoEnergy    DelegateResourceTypo = 1
	DelegateResourceTypoBindWidth DelegateResourceTypo = 0
)

type ExchangeTypo string

const (
	ExchangeTypoUsdt2Trx ExchangeTypo = "USDT2TRX"
	ExchangeTypoTrx2Usdt ExchangeTypo = "TRX2USDT"
)

type OrderStatus string

const (
	OrderStatusCreated      OrderStatus = "已创建"
	OrderStatusNotifyFailed OrderStatus = "通知失败"
	OrderStatusSuccess      OrderStatus = "成功"
)
