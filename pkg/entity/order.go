package entity

type MerchantOrder struct {
	Id              int64  `gorm:"primaryKey;autoIncrement"`
	OrderNo         string `gorm:"size:63;uniqueIndex"`
	MerchantAccount string `gorm:"uniqueIndex:idx_mchaccount_merchorderno_index"`
	MerchantOrderNo string `gorm:"size:63;uniqueIndex:idx_mchaccount_merchorderno_index"`

	LogId         int64  `gorm:""`
	TransactionId string `gorm:"size:255"`
	Chain         string `gorm:"size:15;"`

	Typo         string `gorm:"size:15"`
	Mode         string `gorm:"size:15"`
	Status       string `gorm:"size:15"`
	NotifyStatus string `gorm:"size:15"`
	Currency     string `gorm:"size:15"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	FromAddress string `gorm:"size:255"`
	ToAddress   string `gorm:"size:255"`

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}

type MerchantOrderNotifyLog struct {
	Id              int64 `gorm:"primaryKey;autoIncrement"`
	MerchantOrderId int64 `gorm:""`

	NotifyUrl string `gorm:"size:1023"`

	RequestHeader string
	RequestBody   string

	ResponseHeader string
	ResponseBody   string

	Description string `gorm:"size:2047"`

	TimeAts `gorm:"embedded"`
}
