package converter

import (
	"megichains/pkg/global"

	"github.com/jinzhu/copier"
)

type RespConverter[T any] struct {
	Records []T `json:"records"`
	*PagesBody
}

func ConvertToResp[T any](items []T, current, size int, total int64) (resp *RespConverter[T]) {
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
	MerchOrderId string           `json:"merch_order_id"`
	Chain        global.ChainName `json:"chain"`
	Currency     string           `json:"currency"`
	Receiver     string           `json:"receiver"`
	Seconds      int64            `json:"seconds"`
}

type ChainAddressCreateReq struct {
	Chain string `json:"chain"`
	Count int16  `json:"count"`
}

type OrderListReq struct {
	Pages
	StartEnd
	Id            int64  `form:"id,optional"`
	MerchOrderId  int64  `form:"merch_order_id,optional"`
	TransactionId string `form:"transaction_id,optional"`
	Chain         string `form:"chain,optional"`
	Typo          string `form:"typo,optional"`
	Status        string `form:"status,optional"`
	Currency      string `form:"currency,optional"`
	FromAddress   string `form:"from_address,optional"`
	ToAddress     string `form:"to_address,optional"`
}

type OrderItem struct {
	Id             int64   `json:"id"`
	MerchOrderId   string  `json:"merch_order_id"`
	TransactionId  string  `json:"transaction_id"`
	Chain          string  `json:"chain"`
	Typo           string  `json:"typo"`
	Status         string  `json:"status"`
	Currency       string  `json:"currency"`
	ReceivedAmount float64 `json:"received_amount"`
	ReceivedSun    int64   `json:"received_sun"`
	FromAddress    string  `json:"from_address"`
	ToAddress      string  `json:"to_address"`
	Description    string  `json:"description"`
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
	TimeAts
}

type AddressListReq struct {
	Pages
	StartEnd
	Address string
	Chain   string
	Typo    string
	Status  string
	GroupId int64
}

type AddressListResp struct {
	Records []*AddressItem `json:"records"`
	PagesBody
}
