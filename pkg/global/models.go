// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package global

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

type CollectCalc struct {
	SuccessCount  int64   `gorm:"column:success_count"`
	SuccessAmount float64 `gorm:"column:success_amount"`
}
