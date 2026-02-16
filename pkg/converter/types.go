package converter

import (
	"megichains/pkg/global"

	"github.com/jinzhu/copier"
)

type RespConverter[T any] struct {
	Records []T `json:"records"`
	*PagesBody
}

func ConvertToPagingRecordsResp[T any](items []T, current, size int, total int64) (resp *RespConverter[T]) {
	resp = &RespConverter[T]{
		PagesBody: &PagesBody{
			Current: current,
			Size:    size,
			Total:   total,
		},
		Records: make([]T, 0, size),
	}

	copier.Copy(&resp.Records, &items)

	return
}

func ConvertToRecordsResp[T any](items []T, current, size int, total int64) (resp *RespConverter[T]) {
	resp = &RespConverter[T]{
		Records: make([]T, 0, size),
	}

	copier.Copy(&resp.Records, &items)

	return
}

func ConvertToResp[T any](item T) (resp *T) {
	copier.Copy(&resp, &item)

	return
}

type Pages struct {
	Current int `form:"current"`
	Size    int `form:"size"`
}

type PagesBody struct {
	Current int   `json:"current"`
	Size    int   `json:"size"`
	Total   int64 `json:"total"`
}

type StartEnd struct {
	Start int64 `form:"start,optional"`
	End   int64 `form:"end,optional"`
}

type TimeAts struct {
	UpdatedAt uint64 `json:"update_at"`
	DeletedAt uint64 `json:"delete_at"`
	CreatedAt uint64 `json:"created_at"`
}

type ChainListenReq struct {
	MerchantKey     string
	MerchantOrderNo string           `json:"merchant_order_no"`
	Chain           global.ChainName `json:"chain"`
	Currency        string           `json:"currency"`
	Receiver        string           `json:"receiver"`
	Seconds         int64            `json:"seconds"`
}

type ChainAddressCreateReq struct {
	Chain   string `json:"chain"`
	GroupId int64
	Count   int16 `json:"count"`
}

type OrderListReq struct {
	Pages
	StartEnd
	Id              int64 `form:"id,optional"`
	MerchantKey     string
	MerchantOrderNo int64  `form:"merchant_order_no,optional"`
	TransactionId   string `form:"transaction_id,optional"`
	Chain           string `form:"chain,optional"`
	Typo            string `form:"typo,optional"`
	Status          string `form:"status,optional"`
	Currency        string `form:"currency,optional"`
	FromAddress     string `form:"from_address,optional"`
	ToAddress       string `form:"to_address,optional"`
}

type OrderItem struct {
	Id              int64 `json:"id"`
	MerchantKey     string
	MerchantOrderNo string  `json:"merchant_order_no"`
	TransactionId   string  `json:"transaction_id"`
	Chain           string  `json:"chain"`
	Typo            string  `json:"typo"`
	Status          string  `json:"status"`
	NotifyStatus    string  `json:"notify_status"`
	Currency        string  `json:"currency"`
	ReceivedAmount  float64 `json:"received_amount"`
	ReceivedSun     int64   `json:"received_sun"`
	FromAddress     string  `json:"from_address"`
	ToAddress       string  `json:"to_address"`
	Description     string  `json:"description"`
	TimeAts
}

type OrderListResp struct {
	Records []*OrderItem `json:"records"`
	*PagesBody
}

type TronTransItem struct {
	Id             int64   `json:"id"`
	Chain          string  `json:"chain"`
	Currency       string  `json:"currency"`
	TransactionId  string  `json:"transaction_id"`
	Amount         float64 `json:"amount"`
	Sun            int64   `json:"sun"`
	FromBase58     string  `json:"from_base58"`
	ToBase58       string  `json:"to_base58"`
	Contract       string  `json:"contract"`
	BlockTimestamp int64   `json:"block_timestamp"`
	TimeAts
}

type TronTransListReq struct {
	Pages
	StartEnd
	Id            int64  `form:"id"`
	Currency      string `form:"currency"`
	TransactionId string `form:"transaction_id"`
	FromBase58    string `form:"from_base58"`
	ToBase58      string `form:"to_base58"`
}

type TronTransListResp struct {
	Records []*TronTransItem `json:"records"`
	*PagesBody
}

type EvmLogItem struct {
	Id             int64   `json:"id"`
	Chain          string  `json:"chain"`
	Currency       string  `json:"currency"`
	ChainId        uint64  `json:"chain_id"`
	TxHash         string  `json:"tx_hash"`
	Index          uint    `json:"index"`
	TxIndex        uint    `json:"tx_index"`
	Amount         float64 `json:"amount"`
	Sun            int64   `json:"sun"`
	FromHex        string  `json:"from_hex"`
	ToHex          string  `json:"to_hex"`
	Contract       string  `json:"contract"`
	BlockHash      string  `json:"block_hash"`
	BlockNumber    uint64  `json:"block_number"`
	BlockTimestamp uint64  `json:"block_timestamp"`
	Removed        bool    `json:"removed"`
}

type EvmLogListReq struct {
	Pages
	StartEnd
	Id       int64  `form:"id,optional"`
	Chain    string `json:"chain"`
	Currency string `form:"currency,optional"`
	TxHash   string `form:"tx_hash,optional"`
	FromHex  string `form:"from_hex,optional"`
	ToHex    string `form:"to_hex,optional"`
}

type EvmLogListResp struct {
	Records []*EvmLogItem `json:"records"`
	PagesBody
}

type AddressItem struct {
	Id          int64
	GroupId     int64
	Chain       string
	Typo        string
	Status      string
	Address     string
	Address2    string
	Description string

	AddressBalance

	TimeAts
}

type AddressBalance struct {
	Id        int64
	AddressId int64

	Address string

	TronUsdt   float64
	TronUsdc   float64
	BscUsdt    float64
	BscUsdc    float64
	EthUsdt    float64
	EthUsdc    float64
	SolanaUsdt float64
	SolanaUsdc float64

	TimeAts
}

type AddressWithGroup struct {
	AddressItem
	GroupName string `json:"group_name"`
}

type AddressListReq struct {
	Pages
	StartEnd
	Address  string
	Address2 string
	Chain    string
	Typo     string
	Status   string
	GroupId  int64
}

type AddressListResp struct {
	Records []*AddressItem `json:"records"`
	PagesBody
}

type AddressGroupListReq struct {
	Pages
	Status string `form:"status"`
}

type AddressGroupListResp struct {
	PagesBody
	Records []*AddressGroupItem `json:"records"`
}

type AddressGroupAllResp struct {
	Records []*AddressGroupItem `json:"records"`
}

type AddressGroupItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
	TimeAts
}

type AddressFundListReq struct {
	Pages
	Chain   string
	Address string
}

type AddressFundItem struct {
	Id       int64   `json:"id"`
	Chain    string  `json:"chain"`
	Address  string  `json:"address"`
	TronUsdt float64 `json:"tron_usdt"`
	TronUsdc float64 `json:"tron_usdc"`
	BscUsdt  float64 `json:"bsc_usdt"`
	BscUsdc  float64 `json:"bsc_usdc"`
	EthUsdt  float64 `json:"eth_usdt"`
	EthUsdc  float64 `json:"eth_usdc"`
	TimeAts
}

type AddressFundListResp struct {
	PagesBody
	Records []*AddressFundItem `json:"records"`
}

type AddressFundStatisticsResp struct {
	TronUsdt   float64 `json:"tron_usdt"`
	TronUsdc   float64 `json:"tron_usdc"`
	BscUsdt    float64 `json:"bsc_usdt"`
	BscUsdc    float64 `json:"bsc_usdc"`
	EthUsdt    float64 `json:"eth_usdt"`
	EthUsdc    float64 `json:"eth_usdc"`
	SolanaUsdt float64 `json:"solana_usdt"`
	SolanaUsdc float64 `json:"solana_usdc"`
}

type AddressFundCollectReq struct {
	Chain          string  `json:"chain"`
	Currency       string  `json:"currency"`
	AmountMin      float64 `json:"amount_min"`
	FeeMax         float64 `json:"fee_max"`
	AddressGroupId int64   `json:"address_group_id"`
	SecretKey      string  `json:"secret_key"`
}

type AddressFundCollectResp struct {
	Success bool `json:"success"`
}

type AddressFundCollectListReq struct {
	Pages
	StartEnd
	ReceiverAddress string `form:"receiver_address,optional"`
	AddressGroupId  int64  `form:"address_group_id,optional"`
	Chain           string `form:"chain,optional"`
	Currency        string `form:"currency,optional"`
	Status          string `form:"status,optional"`
}

type AddressFundCollectItem struct {
	Id                  int64   `json:"id"`
	UserId              string  `json:"user_id"`
	Username            string  `json:"username"`
	AddressGroupId      int64   `json:"address_group_id"`
	AddressGroupName    string  `json:"address_group_name"`
	Chain               string  `json:"chain"`
	Currency            string  `json:"currency"`
	Status              string  `json:"status"`
	ReceiverAddress     string  `json:"receiver_address"`
	AmountMin           float64 `json:"amount_min"`
	FeeMax              float64 `json:"fee_max"`
	SuccessAmount       float64 `json:"success_amount"`
	TotalCount          int64   `json:"total_count"`
	SuccessCount        int64   `json:"success_count"`
	TotalGasFee         int64   `json:"total_gas_fee"`
	TotalGasFeeCurrency float64 `json:"total_gas_fee_currency"`
	Description         string  `json:"description"`
	TimeAts
}

type AddressFundCollectListResp struct {
	PagesBody
	Records []*AddressFundCollectItem `json:"records"`
}

type FromAddress struct {
	Address    string
	PrivateKey string
	Balance    float64
}

type AddressFundCollectLogListReq struct {
	Pages
	StartEnd
	CollectId       int64
	FromAddress     string `form:"from_address,optional"`
	ReceiverAddress string `form:"receiver_address,optional"`
	Chain           string `form:"chain,optional"`
	Currency        string `form:"currency,optional"`
	Status          string `form:"status,optional"`
}

type AddressFundCollectLogListResp struct {
	PagesBody
	Records []*AddressFundCollectLogItem `json:"records"`
}

type AddressFundCollectLogItem struct {
	Id                int64
	CollectId         int64
	Chain             string
	Currency          string
	FromAddress       string
	ReceiverAddress   string
	Amount            float64
	Status            string
	TransactionId     string `json:"transaction_id"`
	GasUsed           uint64 `json:"gas_used"`
	GasPrice          int64  `json:"gas_price"`
	EffectiveGasPrice int64  `json:"effective_gas_price"`
	TotalGasFee       int64  `json:"total_gas_fee"`
	Description       string
	TimeAts
}

type UserInfo struct {
	Id       string
	Username string
}

type SolanaTransListReq struct {
	Pages
	StartEnd
	Id            int64  `form:"id"`
	Currency      string `form:"currency"`
	TransactionId string `form:"transaction_id"`
	FromBase58    string `form:"from_base58"`
	ToBase58      string `form:"to_base58"`
}

type SolanaTransListResp struct {
	Records []*SolanaTransItem `json:"records"`
	*PagesBody
}

type SolanaTransItem struct {
	Id            int64   `json:"id"`
	Chain         string  `json:"chain"`
	Currency      string  `json:"currency"`
	TransactionId string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Lamport       int64   `json:"lamport"`
	FromBase58    string  `json:"from_base58"`
	ToBase58      string  `json:"to_base58"`
	Mint          string  `json:"mint"`
	BlockTime     uint64  `json:"block_time"`
}
