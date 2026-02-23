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

type MerchantOrderInteractionLog struct {
	Id              int64 `gorm:"primaryKey;autoIncrement"`
	MerchantOrderId int64 `gorm:""`

	PlaceRequest           string `gorm:"size:8191"`
	PlaceRequestTimestamp  uint64 `gorm:""`
	PlaceResponse          string `gorm:"size:8191"`
	PlaceResponseTimestamp uint64 `gorm:""`

	NotifyRequest           string `gorm:"size:8191"`
	NotifyRequestTimestamp  uint64 `gorm:""`
	NotifyResponse          string `gorm:"size:8191"`
	NotifyResponseTimestamp uint64 `gorm:""`

	Description string `gorm:"size:2047"`

	TimeAts
}
