package global

const (
	StringSuccess = "SUCCESS"
)

type BscTransactionTypo string

const (
	BscTransactionTypoIn  BscTransactionTypo = "入账"
	BscTransactionTypoOut BscTransactionTypo = "出款"
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

type BscTransactionStatus string

const (
	BscTransactionStatusCreated BscTransactionStatus = "已创建"
	BscTransactionStatusSuccess BscTransactionStatus = "成功"
)
