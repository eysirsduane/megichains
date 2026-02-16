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

type OrderNotifyReq struct {
	OrderNo         string  `json:"order_no"`
	MerchantOrderNo string  `json:"merchant_order_no"`
	Status          string  `json:"status"`
	Currency        string  `json:"currency"`
	TxId            string  `json:"txid"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Amount          float64 `json:"amount"`
	Sun             int64   `json:"sun"`
}

type CollectCalc struct {
	SuccessCount  int64   `gorm:"column:success_count"`
	SuccessAmount float64 `gorm:"column:success_amount"`
	TotalGasFee   int64   `gorm:"column:total_gas_fee"`
}
