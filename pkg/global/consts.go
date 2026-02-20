package global

const (
	StringSuccess = "SUCCESS"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
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

type OrderMode string

const (
	OrderModeTest OrderMode = "测试"
	OrderModeProd OrderMode = "正式"
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
	OrderStatusCreated OrderStatus = "CREATED"
	OrderStatusTimeout OrderStatus = "TIMEOUT"
	OrderStatusFailed  OrderStatus = "FAILED"
	OrderStatusSuccess OrderStatus = "SUCCESS"
)

type NotifyStatus string

const (
	NotifyStatusUnknown NotifyStatus = "未知"
	NotifyStatusSuccess NotifyStatus = "成功"
	NotifyStatusFailed  NotifyStatus = "失败"
)

type CollectStatus string

const (
	CollectStatusCreated     CollectStatus = "已创建"
	CollectStatusProcessing  CollectStatus = "处理中"
	CollectStatusPartially   CollectStatus = "部分成功"
	CollectStatusSuccess     CollectStatus = "成功"
	CollectStatusFailed      CollectStatus = "失败"
	CollectLogStatusFinished CollectStatus = "结束"
)

type CollectLogStatus string

const (
	CollectLogStatusCreated    CollectLogStatus = "已创建"
	CollectLogStatusProcessing CollectLogStatus = "处理中"
	CollectLogStatusSuccess    CollectLogStatus = "成功"
	CollectLogStatusTimeout    CollectLogStatus = "超时"
	CollectLogStatusFailed     CollectLogStatus = "失败"
)

type MerchantTypo string

const (
	MerchantTypoTest MerchantTypo = "测试"
	MerchantTypoProd MerchantTypo = "正式"
)
