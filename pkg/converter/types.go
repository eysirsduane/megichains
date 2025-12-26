package converter

import "megichains/pkg/global"

type ChainListenReq struct {
	MerchOrderId string           `json:"merch_order_id"`
	Chain        global.ChainName `json:"chain"`
	Currency     string           `json:"currency"`
	Receiver     string           `json:"receiver"`
	Seconds      int64            `json:"seconds"`
}
